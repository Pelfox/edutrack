package services

import "errors"

var (
	// ErrInvalidInput возвращается при некорректных входных данных.
	ErrInvalidInput = errors.New("invalid input")

	// ErrForbidden возвращается при недостаточных правах для выполнения операции.
	ErrForbidden = errors.New("forbidden")

	// ErrInvalidCredentials возвращается при неверных учётных данных.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUserNotFound возвращается, если пользователь не найден.
	ErrUserNotFound = errors.New("user not found")

	// ErrNotFound возвращается, если запрошенная сущность не найдена.
	ErrNotFound = errors.New("not found")

	// ErrDuplicateUserEmail возвращается при попытке использовать уже занятый email.
	ErrDuplicateUserEmail = errors.New("user email already exists")

	// ErrTokenSigningSecret возвращается, если секрет для подписи JWT-токенов не настроен.
	ErrTokenSigningSecret = errors.New("token signing secret is empty")

	// ErrUnauthenticatedUser возвращается при отсутствии или некорректности авторизации.
	ErrUnauthenticatedUser = errors.New("unauthenticated user")
)
