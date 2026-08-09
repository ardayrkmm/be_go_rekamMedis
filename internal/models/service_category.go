package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type ServiceCategory struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string     `json:"name" firestore:"name"`
	CreatedAt time.Time  `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" firestore:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`
}


func (m *ServiceCategory) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
