package models

import (
	"time"
)

type PasswordResetToken struct {
	Email     string    `json:"email" gorm:"type:varchar(255);index"`
	Token     string    `json:"token" gorm:"primaryKey;type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
}

