package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type TherapySession struct {
	ID                string     `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	AppointmentID     string     `json:"appointment_id" gorm:"type:varchar(36);index"`
	PatientID         string     `json:"patient_id" gorm:"type:varchar(36);index"`
	PhysiotherapistID string     `json:"physiotherapist_id" gorm:"type:varchar(36);index"`
	ServiceMasterID   string     `json:"service_master_id" gorm:"type:varchar(36);index"`
	ServiceMasterIDs  []string   `json:"service_master_ids" gorm:"type:json;serializer:json"`
	TherapyDate       *time.Time `json:"therapy_date"`
	Complaint         string     `json:"complaint"`
	TreatmentGiven    string     `json:"treatment_given"`
	Status            string     `json:"status"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"-" firestore:"DeletedAt"`

	// Relasi — tidak buat FK constraint di DB
	Appointment     *Appointment     `json:"appointment,omitempty" gorm:"-"`
	Patient         *Patient         `json:"patient,omitempty" gorm:"-"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" gorm:"-"`
	ServiceMaster   *ServiceMaster   `json:"service_master,omitempty" gorm:"-"`
	ServiceMasters  []ServiceMaster  `json:"service_masters,omitempty" gorm:"-"`
}

func (m *TherapySession) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
