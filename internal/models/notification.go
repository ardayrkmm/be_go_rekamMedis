package models

import (
	"time"
)

type Notification struct {
	ID              string    `json:"id"` // UUID
	Type            string    `json:"type"`
	NotifiableType  string    `json:"notifiable_type"`
	NotifiableID string `json:"notifiable_id"`
	Data            string    `json:"data"` // MySQL JSON column
	ReadAt          *time.Time`json:"read_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
