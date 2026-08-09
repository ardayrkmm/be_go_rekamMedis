package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type PatientCategory struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


func (m *PatientCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
