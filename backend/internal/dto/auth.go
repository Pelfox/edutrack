package dto

// Login содержит учётные данные для входа в аккаунт.
type Login struct {
	Email    string `json:"email" binding:"required" validate:"required,email"`
	Password string `json:"password" binding:"required" validate:"required,min=8"`
}

// LoginResult содержит результат успешной авторизации пользователя.
type LoginResult struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}
