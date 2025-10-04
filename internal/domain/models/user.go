package models

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Contacts `json:"contacts"`
}

type Contacts struct {
	Email    string  `json:"email"`
	Telegram *string `json:"telegram"`
	ChatId   *int64  `json:"-"`
}
