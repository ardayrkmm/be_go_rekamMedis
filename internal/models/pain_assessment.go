package models

import (
	"github.com/google/uuid"

	"time"

	"gorm.io/gorm"
)

type PainAssessment struct {
	ID string           `gorm:"type:varchar(36);primaryKey" json:"id"`
	TherapySessionID string           `gorm:"type:varchar(36);index" json:"therapy_session_id"`
	PainScale        int            `gorm:"type:int;not null" json:"pain_scale"`
	PainLocation     string         `gorm:"size:255" json:"pain_location"`
	PainType         string         `gorm:"size:255" json:"pain_type"`
	Notes            string         `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	TherapySession *TherapySession `gorm:"foreignKey:TherapySessionID;constraint:-" json:"therapy_session,omitempty"`
}


func (m *PainAssessment) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
