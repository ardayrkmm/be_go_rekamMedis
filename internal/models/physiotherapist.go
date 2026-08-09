package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Physiotherapist struct {
	ID string `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	Name           string         `json:"name"`
	Specialization *string        `json:"specialization"`
	Sip                   string         `json:"sip" gorm:"uniqueIndex;type:varchar(255)"`
	Phone     string         `json:"phone"`
	Email     string         `json:"email"`
	Address   string         `json:"address"`
	Gender    string         `json:"gender"`
	Photo     *string        `json:"photo"`
	Status    string         `json:"status"`
	Password  string         `json:"password,omitempty" gorm:"-" firestore:"-"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Appointments   []Appointment   `json:"appointments,omitempty" gorm:"-"`
	MedicalRecords []MedicalRecord `json:"medical_records,omitempty" gorm:"-"`
}

func (p *Physiotherapist) GetPhotoUrl() *string {
	if p.Photo != nil {
		url := "storage/" + *p.Photo
		return &url
	}
	return nil
}


func (m *Physiotherapist) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
