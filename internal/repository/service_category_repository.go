package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type ServiceCategoryRepository interface {
	FindAll(offset, limit int) ([]models.ServiceCategory, int64, error)
	FindByID(id string) (*models.ServiceCategory, error)
	Create(category *models.ServiceCategory) error
	Update(category *models.ServiceCategory) error
	Delete(id string) error
}

type serviceCategoryRepository struct {
	db *gorm.DB
}

func NewServiceCategoryRepository(db *gorm.DB) ServiceCategoryRepository {
	return &serviceCategoryRepository{db}
}

func (r *serviceCategoryRepository) FindAll(offset, limit int) ([]models.ServiceCategory, int64, error) {
	var categories []models.ServiceCategory
	var total int64

	err := r.db.Model(&models.ServiceCategory{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&categories).Error
	return categories, total, err
}

func (r *serviceCategoryRepository) FindByID(id string) (*models.ServiceCategory, error) {
	var category models.ServiceCategory
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *serviceCategoryRepository) Create(item *models.ServiceCategory) error {
	return r.db.Create(item).Error
}

func (r *serviceCategoryRepository) Update(item *models.ServiceCategory) error {
	return r.db.Save(item).Error
}

func (r *serviceCategoryRepository) Delete(id string) error {
	return r.db.Delete(&models.ServiceCategory{}, "id = ?", id).Error
}
