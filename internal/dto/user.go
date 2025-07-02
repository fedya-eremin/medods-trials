package dto

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id" db:"id"`
}
