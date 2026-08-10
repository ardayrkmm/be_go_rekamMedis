package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type PainAssessmentRepository interface {
	FindAll(offset, limit int) ([]models.PainAssessment, int64, error)
	FindByID(id string) (*models.PainAssessment, error)
	Create(assessment *models.PainAssessment) error
	Update(assessment *models.PainAssessment) error
	Delete(id string) error
}

type painAssessmentRepository struct {
	db *gorm.DB
}

func NewPainAssessmentRepository(db *gorm.DB) PainAssessmentRepository {
	return &painAssessmentRepository{db}
}

func (r *painAssessmentRepository) FindAll(offset, limit int) ([]models.PainAssessment, int64, error) {
	var assessments []models.PainAssessment
	var total int64
	err := r.db.Model(&models.PainAssessment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Preload("TherapySession").Offset(offset).Limit(limit).Find(&assessments).Error
	return assessments, total, err
}

func (r *painAssessmentRepository) FindByID(id string) (*models.PainAssessment, error) {
	var assessment models.PainAssessment
	err := r.db.Preload("TherapySession").First(&assessment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &assessment, nil
}

func (r *painAssessmentRepository) Create(assessment *models.PainAssessment) error {
	return r.db.Create(assessment).Error
}

func (r *painAssessmentRepository) Update(assessment *models.PainAssessment) error {
	return r.db.Save(assessment).Error
}

func (r *painAssessmentRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.PainAssessment{}).Error
}

