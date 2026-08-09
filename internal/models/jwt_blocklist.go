package models

import (
	"time"
)

type JwtBlocklist struct {
	Token     string    `json:"token" gorm:"primaryKey;type:varchar(500)"`
	ExpiresAt time.Time `json:"expires_at"`
}

