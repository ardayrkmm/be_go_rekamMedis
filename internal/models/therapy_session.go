package models

import (
	"time"
)

type TherapySession struct {
	ID string `firestore:"id,omitempty" json:"id"`
	AppointmentID      string           `json:"appointment_id"`
	PatientID string `json:"patient_id"`
	PhysiotherapistID string `json:"physiotherapist_id"`
	ServiceMasterID    string           `json:"service_master_id"`
	ServiceMasterIDs   []string         `json:"service_master_ids,omitempty" firestore:"service_master_ids,omitempty"`
	TherapyDate        *time.Time       `json:"therapy_date"`
	Complaint          string           `json:"complaint"`
	TreatmentGiven     string           `json:"treatment_given"`
	Status             string           `json:"status"`
	Notes              string           `json:"notes"`
	CreatedAt          time.Time        `json:"created_at"`
	UpdatedAt          time.Time        `json:"updated_at"`
	DeletedAt          *time.Time       `json:"-" firestore:"DeletedAt"`

	Appointment        *Appointment       `json:"appointment,omitempty"`
	Patient            *Patient           `json:"patient,omitempty"`
	Physiotherapist    *Physiotherapist   `json:"physiotherapist,omitempty"`
	ServiceMaster      *ServiceMaster     `json:"service_master,omitempty"`
	ServiceMasters     []ServiceMaster    `json:"service_masters,omitempty"`
}
