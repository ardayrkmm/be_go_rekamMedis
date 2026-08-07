package models

import (
	"time"
)

type Physiotherapist struct {
	ID string `firestore:"id,omitempty" json:"id"`
	Name           string         `json:"name"`
	Specialization *string        `json:"specialization"`
	Sip            string         `json:"sip"`
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

	Appointments   []Appointment   `json:"appointments,omitempty"`
	MedicalRecords []MedicalRecord `json:"medical_records,omitempty"`
}

func (p *Physiotherapist) GetPhotoUrl() *string {
	if p.Photo != nil {
		url := "storage/" + *p.Photo
		return &url
	}
	return nil
}
