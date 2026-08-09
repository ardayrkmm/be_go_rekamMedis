package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type ActivityLog struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	LogName     string    `json:"log_name"`
	Description string    `json:"description"`
	SubjectType string    `json:"subject_type"`
	SubjectID   *int     `json:"subject_id"`
	CauserType  string    `json:"causer_type"`
	CauserID    *int     `json:"causer_id"`
	Properties  string    `json:"properties"` // MySQL JSON column
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}


func (m *ActivityLog) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
