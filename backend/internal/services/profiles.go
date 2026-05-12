package services

import (
	"context"

	"github.com/Pelfox/edutrack/backend/internal/dto"
	"github.com/Pelfox/edutrack/backend/internal/repositories"
)

// ProfileService описывает операции с профилем текущего пользователя.
type ProfileService interface {
	// GetMe возвращает профиль текущего пользователя.
	GetMe(ctx context.Context, actor dto.Actor) (*dto.Profile, error)
}

type profileService struct {
	profiles repositories.ProfileRepository
}

// NewProfileService создаёт сервис профилей пользователей.
func NewProfileService(profiles repositories.ProfileRepository) ProfileService {
	return &profileService{profiles: profiles}
}

func (service *profileService) GetMe(ctx context.Context, actor dto.Actor) (*dto.Profile, error) {
	if actor.ID == "" {
		return nil, ErrUnauthenticatedUser
	}

	profile, err := service.profiles.GetByUserID(ctx, actor.ID, actor.Role)
	if err != nil {
		return nil, mapDirectoryRepositoryError(err)
	}

	return toProfileOutput(profile), nil
}

func toProfileOutput(profile *repositories.Profile) *dto.Profile {
	return &dto.Profile{
		ID:         profile.ID,
		UserID:     profile.UserID,
		Email:      profile.Email,
		Role:       profile.Role,
		LastName:   profile.LastName,
		FirstName:  profile.FirstName,
		MiddleName: profile.MiddleName,
	}
}
