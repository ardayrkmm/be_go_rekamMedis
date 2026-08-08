package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
)

type ServiceCategoryUseCase interface {
	Fetch(offset, limit int) ([]models.ServiceCategory, int64, error)
	GetByID(id string) (*models.ServiceCategory, error)
	Store(category *models.ServiceCategory) error
	Update(id string, category *models.ServiceCategory) error
	Delete(id string) error
}

type serviceCategoryUseCase struct {
	categoryRepo repository.ServiceCategoryRepository
}

func NewServiceCategoryUseCase(categoryRepo repository.ServiceCategoryRepository) ServiceCategoryUseCase {
	return &serviceCategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

func (u *serviceCategoryUseCase) Fetch(offset, limit int) ([]models.ServiceCategory, int64, error) {
	return u.categoryRepo.FindAll(offset, limit)
}

func (u *serviceCategoryUseCase) GetByID(id string) (*models.ServiceCategory, error) {
	return u.categoryRepo.FindByID(id)
}

func (u *serviceCategoryUseCase) Store(category *models.ServiceCategory) error {
	return u.categoryRepo.Create(category)
}

func (u *serviceCategoryUseCase) Update(id string, req *models.ServiceCategory) error {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return err
	}

	category.Name = req.Name
	return u.categoryRepo.Update(category)
}

func (u *serviceCategoryUseCase) Delete(id string) error {
	return u.categoryRepo.Delete(id)
}
