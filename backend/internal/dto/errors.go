package dto

// Error содержит унифицированный формат HTTP-ошибки.
type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
