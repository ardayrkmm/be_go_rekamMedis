package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Appointment struct {
	ID                string     `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	VisitNumber       *string    `json:"visit_number"`
	PatientID         string     `json:"patient_id" gorm:"type:varchar(36);index"`
	PhysiotherapistID string     `json:"physiotherapist_id" gorm:"type:varchar(36);index"`
	ServiceMasterID   string     `json:"service_master_id" gorm:"type:varchar(36);index"`
	AppointmentDate   *time.Time `json:"appointment_date"`
	AppointmentTime   string     `json:"appointment_time"`
	Complaint         string     `json:"complaint"`
	Status            string     `json:"status"`
	Notes             string     `json:"notes"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"-" firestore:"DeletedAt"`

	// Relasi — tidak buat FK constraint di DB
	Patient         *Patient         `json:"patient,omitempty" gorm:"-"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" gorm:"-"`
	ServiceMaster   *ServiceMaster   `json:"service_master,omitempty" gorm:"-"`
	TherapySession  *TherapySession  `json:"therapy_session,omitempty" gorm:"-"`
}

func (m *Appointment) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
