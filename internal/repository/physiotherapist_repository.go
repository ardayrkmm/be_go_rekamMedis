package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type PhysiotherapistRepository interface {
	FindAll(offset, limit int) ([]models.Physiotherapist, int64, error)
	FindByID(id string) (*models.Physiotherapist, error)
	Create(physio *models.Physiotherapist) error
	Update(physio *models.Physiotherapist) error
	Delete(id string) error
	Restore(id string) error
}

type physiotherapistRepository struct {
	db *gorm.DB
}

func NewPhysiotherapistRepository(db *gorm.DB) PhysiotherapistRepository {
	return &physiotherapistRepository{db}
}

func (r *physiotherapistRepository) FindAll(offset, limit int) ([]models.Physiotherapist, int64, error) {
	var physios []models.Physiotherapist
	var total int64

	err := r.db.Model(&models.Physiotherapist{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&physios).Error
	return physios, total, err
}

func (r *physiotherapistRepository) FindByID(id string) (*models.Physiotherapist, error) {
	var physio models.Physiotherapist
	err := r.db.First(&physio, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &physio, nil
}

func (r *physiotherapistRepository) Create(physio *models.Physiotherapist) error {
	return r.db.Create(physio).Error
}

func (r *physiotherapistRepository) Update(physio *models.Physiotherapist) error {
	return r.db.Save(physio).Error
}

func (r *physiotherapistRepository) Delete(id string) error {
	return r.db.Delete(&models.Physiotherapist{}, id).Error
}

func (r *physiotherapistRepository) Restore(id string) error {
	return r.db.Unscoped().Model(&models.Physiotherapist{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
