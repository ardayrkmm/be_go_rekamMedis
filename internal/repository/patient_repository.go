package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type PatientRepository interface {
	FindAll(offset, limit int) ([]models.Patient, int64, error)
	FindByID(id string) (*models.Patient, error)
	Create(patient *models.Patient) error
	Update(patient *models.Patient) error
	Delete(id string) error
	Restore(id string) error
}

type patientRepository struct {
	db *gorm.DB
}

func NewPatientRepository(db *gorm.DB) PatientRepository {
	return &patientRepository{db}
}

func (r *patientRepository) FindAll(offset, limit int) ([]models.Patient, int64, error) {
	var patients []models.Patient
	var total int64

	err := r.db.Model(&models.Patient{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Category").Preload("GenderData").Offset(offset).Limit(limit).Find(&patients).Error
	return patients, total, err
}

func (r *patientRepository) FindByID(id string) (*models.Patient, error) {
	var patient models.Patient
	err := r.db.Preload("Category").Preload("GenderData").First(&patient, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

func (r *patientRepository) Create(patient *models.Patient) error {
	return r.db.Create(patient).Error
}

func (r *patientRepository) Update(patient *models.Patient) error {
	return r.db.Save(patient).Error
}

func (r *patientRepository) Delete(id string) error {
	return r.db.Delete(&models.Patient{}, id).Error
}

func (r *patientRepository) Restore(id string) error {
	return r.db.Unscoped().Model(&models.Patient{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
