package models

import (
	"time"
)

type Payment struct {
	ID string `firestore:"id,omitempty" json:"id"`
	InvoiceNumber      string          `json:"invoice_number"`
	AppointmentID    string `json:"appointment_id,omitempty" firestore:"appointment_id"`
	TherapySessionID string `json:"therapy_session_id"`
	PatientID string `json:"patient_id" firestore:"patient_id"`
	PatientName string `json:"patient_name" firestore:"patient_name"`
	PhysiotherapistID string `json:"physiotherapist_id" firestore:"physiotherapist_id"`
	PhysiotherapistName string `json:"physiotherapist_name" firestore:"physiotherapist_name"`
	PaymentDate        *time.Time      `json:"payment_date"`
	PaymentMethod      string          `json:"payment_method"`
	Status             string          `json:"status"`
	Subtotal           float64         `json:"subtotal"`
	Discount           float64         `json:"discount"`
	Tax                float64         `json:"tax"`
	Total              float64         `json:"total"`
	Notes              string          `json:"notes"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	TherapySession  *TherapySession  `json:"therapy_session,omitempty"`
	Appointment     *Appointment     `json:"appointment,omitempty"`
	Patient         *Patient         `json:"patient,omitempty"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty"`
	PaymentDetails  []PaymentDetail  `json:"payment_details,omitempty"`
}
