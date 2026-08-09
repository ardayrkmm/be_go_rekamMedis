package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type MedicalRecord struct {
	ID                string     `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	VisitNumber       *string    `json:"visit_number"`
	PatientID         string     `json:"patient_id" gorm:"type:varchar(36);index"`
	ServiceID         *string    `json:"service_id"`
	PhysiotherapistID string     `json:"physiotherapist_id" gorm:"type:varchar(36);index"`
	AppointmentID     *string    `json:"appointment_id"`
	ExaminationDate   *time.Time `json:"examination_date"`
	Anamnesis         string     `json:"anamnesis"`
	Diagnosis         string     `json:"diagnosis"`
	Therapy           string     `json:"therapy"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"-" firestore:"DeletedAt"`

	// Relasi — tidak buat FK constraint di DB
	Patient         *Patient         `json:"patient,omitempty" gorm:"-"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" gorm:"-"`
	Service         *ServiceMaster   `json:"service,omitempty" gorm:"-"`
}

func (m *MedicalRecord) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
