package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Notification struct {
	ID string `json:"id" gorm:"type:varchar(36);primaryKey"` // UUID
	Type            string    `json:"type"`
	NotifiableType  string    `json:"notifiable_type"`
	NotifiableID string `json:"notifiable_id" gorm:"type:varchar(36)"`
	Data            string    `json:"data"` // MySQL JSON column
	ReadAt          *time.Time`json:"read_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}


func (m *Notification) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
