package models

import (
	"time"
)

type JwtBlocklist struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
