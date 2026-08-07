package models

import (
	"time"
)

type MedicalRecord struct {
	ID string `firestore:"id,omitempty" json:"id"`
	VisitNumber       *string        `json:"visit_number"`
	PatientID string `json:"patient_id"`
	ServiceID         *string          `json:"service_id"`
	PhysiotherapistID string `json:"physiotherapist_id"`
	AppointmentID     *string          `json:"appointment_id"`
	ExaminationDate   *time.Time     `json:"examination_date"`
	Anamnesis         string         `json:"anamnesis"`
	Diagnosis         string         `json:"diagnosis"`
	Therapy           string         `json:"therapy"`
	Notes             string         `json:"notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Patient         *Patient         `json:"patient,omitempty" firestore:"-"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" firestore:"-"`
	Service         *ServiceMaster   `json:"service,omitempty" firestore:"-"`
}
