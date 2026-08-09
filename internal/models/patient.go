package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"time"
)

type Patient struct {
	ID                    string     `firestore:"id,omitempty" json:"id" gorm:"type:varchar(36);primaryKey"`
	MedicalRecordNumber   string     `json:"medical_record_number" gorm:"uniqueIndex;type:varchar(255)"`
	Nik                   *string     `json:"nik" gorm:"uniqueIndex;type:varchar(255)"`
	Name                  string     `json:"name" gorm:"type:varchar(255)"`
	BirthDate             *time.Time `json:"birth_date"`
	PatientCategoryID     string     `json:"patient_category_id" gorm:"type:varchar(36);index"`
	GenderID              string     `json:"gender_id" gorm:"type:varchar(36);index"`
	BloodType             string     `json:"blood_type"`
	Address               string     `json:"address"`
	Phone                 string     `json:"phone"`
	Email                 *string     `json:"email" gorm:"uniqueIndex;type:varchar(255)"`
	Occupation            string     `json:"occupation"`
	MaritalStatus         string     `json:"marital_status"`
	EmergencyContactName  string     `json:"emergency_contact_name"`
	EmergencyContactPhone string     `json:"emergency_contact_phone"`
	MedicalHistory        string     `json:"medical_history"`
	Allergies             string     `json:"allergies"`
	CreatedAt             time.Time  `json:"created_at"`
	UpdatedAt             time.Time  `json:"updated_at"`
	DeletedAt             *time.Time `json:"-" firestore:"DeletedAt"`

	// Relasi
	Category       *PatientCategory `json:"category,omitempty" gorm:"foreignKey:PatientCategoryID;constraint:-"`
	GenderData     *Gender          `json:"gender_data,omitempty" gorm:"foreignKey:GenderID;constraint:-"`
	Appointments []Appointment `json:"appointments,omitempty" gorm:"foreignKey:PatientID;constraint:-"`
	MedicalRecords []MedicalRecord `json:"medical_records,omitempty" gorm:"foreignKey:PatientID;constraint:-"`
}

func (m *Patient) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Nik != nil && *m.Nik == "" { m.Nik = nil }
	if m.Email != nil && *m.Email == "" { m.Email = nil }
	return
}



func (m *Patient) BeforeUpdate(tx *gorm.DB) (err error) {
	if m.Nik != nil && *m.Nik == "" { m.Nik = nil }
	if m.Email != nil && *m.Email == "" { m.Email = nil }
	return
}
