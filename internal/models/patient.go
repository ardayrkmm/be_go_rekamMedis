package models

import (
	"time"
)

type Patient struct {
	ID string `firestore:"id,omitempty" json:"id"`
	MedicalRecordNumber   string         `json:"medical_record_number"`
	Nik                   string         `json:"nik"`
	Name                  string         `json:"name"`
	BirthDate             *time.Time     `json:"birth_date"`
	PatientCategoryID     string          `json:"patient_category_id"`
	Category              *PatientCategory `json:"category,omitempty"`
	GenderID              string          `json:"gender_id"`
	GenderData            *Gender        `json:"gender_data,omitempty"`
	BloodType             string         `json:"blood_type"`
	Address               string         `json:"address"`
	Phone                 string         `json:"phone"`
	Email                 string         `json:"email"`
	Occupation            string         `json:"occupation"`
	MaritalStatus         string         `json:"marital_status"`
	EmergencyContactName  string         `json:"emergency_contact_name"`
	EmergencyContactPhone string         `json:"emergency_contact_phone"`
	MedicalHistory        string         `json:"medical_history"`
	Allergies             string         `json:"allergies"`
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt *time.Time `json:"-" firestore:"DeletedAt"`

	Appointments   []Appointment   `json:"appointments,omitempty"`
	MedicalRecords []MedicalRecord `json:"medical_records,omitempty"`
}
