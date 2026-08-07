package models

import (
	"time"
)

type Appointment struct {
	ID string `firestore:"id,omitempty" json:"id"`
	PatientID string `json:"patient_id"`
	PhysiotherapistID string `json:"physiotherapist_id"`
	ServiceMasterID    string         `json:"service_master_id"`
	AppointmentDate    *time.Time      `json:"appointment_date"`
	AppointmentTime    string          `json:"appointment_time"`
	Complaint          string          `json:"complaint"`
	Status             string          `json:"status"`
	Notes              string          `json:"notes"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Patient         *Patient         `json:"patient,omitempty"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty"`
	ServiceMaster   *ServiceMaster   `json:"service_master,omitempty"`
}
