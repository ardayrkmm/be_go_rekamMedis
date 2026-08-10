package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"time"
)

type MasterUseCase interface {
	GetPatientCategories() ([]models.PatientCategory, error)
	GetGenders() ([]models.Gender, error)
	CreatePatientCategory(category *models.PatientCategory) error
	UpdatePatientCategory(id string, category *models.PatientCategory) error
	DeletePatientCategory(id string) error
}

type masterUseCase struct {
	masterRepo repository.MasterRepository
}

func NewMasterUseCase(masterRepo repository.MasterRepository) MasterUseCase {
	return &masterUseCase{masterRepo}
}

func (u *masterUseCase) GetPatientCategories() ([]models.PatientCategory, error) {
	return u.masterRepo.GetPatientCategories()
}

func (u *masterUseCase) CreatePatientCategory(category *models.PatientCategory) error {
	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()
	return u.masterRepo.CreatePatientCategory(category)
}

func (u *masterUseCase) UpdatePatientCategory(id string, category *models.PatientCategory) error {
	category.UpdatedAt = time.Now()
	return u.masterRepo.UpdatePatientCategory(id, category)
}

func (u *masterUseCase) DeletePatientCategory(id string) error {
	return u.masterRepo.DeletePatientCategory(id)
}

func (u *masterUseCase) GetGenders() ([]models.Gender, error) {
	return u.masterRepo.GetGenders()
}
