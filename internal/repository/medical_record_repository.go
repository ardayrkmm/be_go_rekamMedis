package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type MedicalRecordRepository interface {
	FindAll(offset, limit int, search string) ([]models.MedicalRecord, int64, error)
	FindByID(id string) (*models.MedicalRecord, error)
	FindByPatientID(patientID string) ([]models.MedicalRecord, error)
	Create(record *models.MedicalRecord) error
	Update(record *models.MedicalRecord) error
	Delete(id string) error
}

type medicalRecordRepository struct {
	db *gorm.DB
}

func NewMedicalRecordRepository(db *gorm.DB) MedicalRecordRepository {
	return &medicalRecordRepository{db}
}

func (r *medicalRecordRepository) FindAll(offset, limit int, search string) ([]models.MedicalRecord, int64, error) {
	var records []models.MedicalRecord
	var total int64

	query := r.db.Model(&models.MedicalRecord{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("LEFT JOIN patients ON patients.id = medical_records.patient_id").
			Where("medical_records.visit_number LIKE ? OR medical_records.diagnosis LIKE ? OR medical_records.anamnesis LIKE ? OR patients.name LIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Patient").Preload("Patient.GenderData").Preload("Patient.GenderData").Preload("Physiotherapist").Preload("Service").Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

func (r *medicalRecordRepository) FindByID(id string) (*models.MedicalRecord, error) {
	var record models.MedicalRecord
	err := r.db.Preload("Patient").Preload("Patient.GenderData").Preload("Patient.GenderData").Preload("Physiotherapist").Preload("Service").First(&record, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *medicalRecordRepository) FindByPatientID(patientID string) ([]models.MedicalRecord, error) {
	var records []models.MedicalRecord
	err := r.db.Where("patient_id = ?", patientID).Preload("Physiotherapist").Preload("Service").Find(&records).Error
	return records, err
}

func (r *medicalRecordRepository) Create(record *models.MedicalRecord) error {
	return r.db.Create(record).Error
}

func (r *medicalRecordRepository) Update(record *models.MedicalRecord) error {
	return r.db.Save(record).Error
}

func (r *medicalRecordRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.MedicalRecord{}).Error
}

