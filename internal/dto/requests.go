package dto

type CreateUserRequest struct {
	Name     string  `json:"name" validate:"required,min=1"`
	Email    string  `json:"email" validate:"required,email"`
	Telegram *string `json:"telegram" validate:"omitempty,min=1"`
}

type UpdateIPhoneRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	NewPrice int    `json:"new_price" validate:"required,min=1"`
}
