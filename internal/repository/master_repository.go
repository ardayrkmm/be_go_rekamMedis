package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type MasterRepository interface {
	GetPatientCategories() ([]models.PatientCategory, error)
	GetGenders() ([]models.Gender, error)
	CreatePatientCategory(category *models.PatientCategory) error
	UpdatePatientCategory(id string, category *models.PatientCategory) error
	DeletePatientCategory(id string) error
	CreateGender(gender *models.Gender) error
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db}
}

func (r *masterRepository) GetPatientCategories() ([]models.PatientCategory, error) {
	var categories []models.PatientCategory
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *masterRepository) GetGenders() ([]models.Gender, error) {
	var genders []models.Gender
	err := r.db.Find(&genders).Error
	return genders, err
}

func (r *masterRepository) CreatePatientCategory(category *models.PatientCategory) error {
	return r.db.Create(category).Error
}

func (r *masterRepository) UpdatePatientCategory(id string, category *models.PatientCategory) error {
	return r.db.Model(&models.PatientCategory{}).Where("id = ?", id).Updates(category).Error
}

func (r *masterRepository) DeletePatientCategory(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.PatientCategory{}).Error
}

func (r *masterRepository) CreateGender(gender *models.Gender) error {
	return r.db.Create(gender).Error
}
