package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
)

type UserUseCase interface {
	Fetch(offset, limit int) ([]models.User, int64, error)
	GetByID(id string) (*models.User, error)
	Store(user *models.User) error
	Update(id string, req *models.User) error
	Delete(id string) error
	Restore(id string) error
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{
		userRepo: userRepo,
	}
}

func (u *userUseCase) Fetch(offset, limit int) ([]models.User, int64, error) {
	return u.userRepo.FindAll(offset, limit)
}

func (u *userUseCase) GetByID(id string) (*models.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUseCase) Store(user *models.User) error {
	return u.userRepo.Create(user)
}

func (u *userUseCase) Update(id string, req *models.User) error {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Photo != nil {
		user.Photo = req.Photo
	}

	return u.userRepo.Update(user)
}

func (u *userUseCase) Delete(id string) error {
	return u.userRepo.Delete(id)
}

func (u *userUseCase) Restore(id string) error {
	return u.userRepo.Restore(id)
}
