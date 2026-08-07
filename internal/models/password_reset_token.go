package models

import (
	"time"
)

type PasswordResetToken struct {
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
