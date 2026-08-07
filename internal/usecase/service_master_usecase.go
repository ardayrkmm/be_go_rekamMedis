package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"time"
)

type ServiceMasterUseCase interface {
	Fetch(offset, limit int) ([]models.ServiceMaster, int64, error)
	GetByID(id string) (*models.ServiceMaster, error)
	Store(service *models.ServiceMaster) error
	Update(id string, service *models.ServiceMaster) error
	Delete(id string) error
}

type serviceMasterUseCase struct {
	serviceRepo repository.ServiceMasterRepository
}

func NewServiceMasterUseCase(serviceRepo repository.ServiceMasterRepository) ServiceMasterUseCase {
	return &serviceMasterUseCase{
		serviceRepo: serviceRepo,
	}
}

func (u *serviceMasterUseCase) Fetch(offset, limit int) ([]models.ServiceMaster, int64, error) {
	return u.serviceRepo.FindAll(offset, limit)
}

func (u *serviceMasterUseCase) GetByID(id string) (*models.ServiceMaster, error) {
	return u.serviceRepo.FindByID(id)
}

func (u *serviceMasterUseCase) Store(service *models.ServiceMaster) error {
	if service.Code == "" {
		service.Code = "SVC-" + time.Now().Format("20060102150405")
	}
	return u.serviceRepo.Create(service)
}

func (u *serviceMasterUseCase) Update(id string, req *models.ServiceMaster) error {
	service, err := u.serviceRepo.FindByID(id)
	if err != nil {
		return err
	}

	service.Name = req.Name
	service.Description = req.Description
	service.BasePrice = req.BasePrice
	service.IsActive = req.IsActive
	service.Category = req.Category
	service.Duration = req.Duration
	if req.Code != "" {
		service.Code = req.Code
	} else if service.Code == "" {
		service.Code = "SVC-" + time.Now().Format("20060102150405")
	}

	return u.serviceRepo.Update(service)
}

func (u *serviceMasterUseCase) Delete(id string) error {
	return u.serviceRepo.Delete(id)
}
