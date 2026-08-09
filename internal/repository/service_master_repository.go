package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type ServiceMasterRepository interface {
	FindAll(offset, limit int) ([]models.ServiceMaster, int64, error)
	FindByID(id string) (*models.ServiceMaster, error)
	Create(service *models.ServiceMaster) error
	Update(service *models.ServiceMaster) error
	Delete(id string) error
}

type serviceMasterRepository struct {
	db *gorm.DB
}

func NewServiceMasterRepository(db *gorm.DB) ServiceMasterRepository {
	return &serviceMasterRepository{db}
}

func (r *serviceMasterRepository) FindAll(offset, limit int) ([]models.ServiceMaster, int64, error) {
	var services []models.ServiceMaster
	var total int64

	err := r.db.Model(&models.ServiceMaster{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&services).Error
	return services, total, err
}

func (r *serviceMasterRepository) FindByID(id string) (*models.ServiceMaster, error) {
	var service models.ServiceMaster
	err := r.db.First(&service, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &service, nil
}

func (r *serviceMasterRepository) Create(service *models.ServiceMaster) error {
	return r.db.Create(service).Error
}

func (r *serviceMasterRepository) Update(service *models.ServiceMaster) error {
	return r.db.Save(service).Error
}

func (r *serviceMasterRepository) Delete(id string) error {
	return r.db.Delete(&models.ServiceMaster{}, id).Error
}
