package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// AuthService описывает операции авторизации и проверки JWT-токенов.
type AuthService interface {
	// Login проверяет учётные данные пользователя и возвращает JWT-токен.
	Login(ctx context.Context, input dto.Login) (*dto.LoginResult, error)

	// ParseToken проверяет JWT-токен и возвращает данные авторизованного пользователя.
	ParseToken(tokenValue string) (dto.Actor, error)
}

type authService struct {
	users     repositories.AuthRepository
	jwtSecret string
	validator *validator.Validate
	logger    zerolog.Logger
}

// Claims содержит пользовательские и стандартные поля JWT-токена.
type Claims struct {
	UserID string                `json:"user_id"`
	Role   repositories.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// NewAuthService создаёт сервис авторизации.
func NewAuthService(users repositories.AuthRepository, jwtSecret string, logger zerolog.Logger) AuthService {
	return &authService{
		users:     users,
		jwtSecret: jwtSecret,
		validator: validator.New(validator.WithRequiredStructEnabled()),
		logger:    logger,
	}
}

func (service *authService) Login(ctx context.Context, input dto.Login) (*dto.LoginResult, error) {
	input.Email = normalizeEmail(input.Email)
	if err := validateStruct(service.validator, input); err != nil {
		service.logger.Warn().Err(err).Str("email", input.Email).Msg("login validation failed")
		return nil, err
	}

	user, err := service.users.GetByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			service.logger.Warn().Str("email", input.Email).Msg("login failed")
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		service.logger.Warn().
			Str("user_id", user.ID).
			Str("email", user.Email).
			Msg("login failed")
		return nil, ErrInvalidCredentials
	}

	token, err := service.createToken(user.ID, user.Role)
	if err != nil {
		service.logger.Error().Err(err).Str("user_id", user.ID).Msg("failed to create auth token")
		return nil, err
	}

	service.logger.Info().
		Str("user_id", user.ID).
		Str("email", user.Email).
		Str("role", string(user.Role)).
		Msg("user logged in")

	return &dto.LoginResult{
		User:  *toUserOutput(user),
		Token: token,
	}, nil
}

func (service *authService) createToken(userID string, role repositories.UserRole) (string, error) {
	if service.jwtSecret == "" {
		return "", ErrTokenSigningSecret
	}

	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)),
		},
	})

	signedToken, err := token.SignedString([]byte(service.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (service *authService) ParseToken(tokenValue string) (dto.Actor, error) {
	if service.jwtSecret == "" {
		return dto.Actor{}, ErrTokenSigningSecret
	}

	token, err := jwt.ParseWithClaims(tokenValue, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnauthenticatedUser
		}

		return []byte(service.jwtSecret), nil
	})
	if err != nil {
		return dto.Actor{}, ErrUnauthenticatedUser
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return dto.Actor{}, ErrUnauthenticatedUser
	}
	if err := validateUUID(claims.UserID); err != nil {
		return dto.Actor{}, ErrUnauthenticatedUser
	}
	if err := service.validator.Var(claims.Role, "required,oneof=administrator teacher student"); err != nil {
		return dto.Actor{}, ErrUnauthenticatedUser
	}

	return dto.Actor{ID: claims.UserID, Role: claims.Role}, nil
}
