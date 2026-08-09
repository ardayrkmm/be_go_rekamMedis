package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Payment struct {
	ID                  string     `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	InvoiceNumber       string     `json:"invoice_number"`
	AppointmentID       string     `json:"appointment_id,omitempty" firestore:"appointment_id" gorm:"type:varchar(36);index"`
	TherapySessionID    string     `json:"therapy_session_id" gorm:"type:varchar(36);index"`
	PatientID           string     `json:"patient_id" firestore:"patient_id" gorm:"type:varchar(36);index"`
	PatientName         string     `json:"patient_name" firestore:"patient_name"`
	PhysiotherapistID   string     `json:"physiotherapist_id" firestore:"physiotherapist_id" gorm:"type:varchar(36);index"`
	PhysiotherapistName string     `json:"physiotherapist_name" firestore:"physiotherapist_name"`
	PaymentDate         *time.Time `json:"payment_date"`
	PaymentMethod       string     `json:"payment_method"`
	Status              string     `json:"status"`
	Subtotal            float64    `json:"subtotal"`
	Discount            float64    `json:"discount"`
	Tax                 float64    `json:"tax"`
	Total               float64    `json:"total"`
	Notes               string     `json:"notes"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	DeletedAt           *time.Time `json:"-" firestore:"DeletedAt"`

	// Relasi — tidak buat FK constraint di DB
	TherapySession *TherapySession `json:"therapy_session,omitempty" gorm:"foreignKey:AppointmentID;constraint:-"`
	Appointment *Appointment `json:"appointment,omitempty" gorm:"foreignKey:AppointmentID;constraint:-"`
	Patient *Patient `json:"patient,omitempty" gorm:"foreignKey:PatientID;constraint:-"`
	Physiotherapist *Physiotherapist `json:"physiotherapist,omitempty" gorm:"foreignKey:PhysiotherapistID;constraint:-"`
	PaymentDetails []PaymentDetail `json:"payment_details,omitempty" gorm:"foreignKey:PaymentID;constraint:-"`
}

func (m *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return
}
