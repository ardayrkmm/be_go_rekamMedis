package models

import (
	"github.com/google/uuid"

	"time"

	"gorm.io/gorm"
)

type ExerciseProgram struct {
	ID string           `gorm:"type:varchar(36);primaryKey" json:"id"`
	TherapySessionID string           `gorm:"type:varchar(36);index" json:"therapy_session_id"`
	ExerciseName     string         `gorm:"size:255;not null" json:"exercise_name"`
	Description      string         `gorm:"type:text" json:"description"`
	Repetitions      int            `gorm:"default:0" json:"repetitions"`
	Sets             int            `gorm:"default:0" json:"sets"`
	Frequency        string         `gorm:"size:100" json:"frequency"`
	Notes            string         `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	TherapySession *TherapySession `gorm:"foreignKey:TherapySessionID;constraint:-" json:"therapy_session,omitempty"`
}


func (m *ExerciseProgram) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
