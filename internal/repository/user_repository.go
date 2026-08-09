package repository

import (
	"backend_go/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(offset, limit int) ([]models.User, int64, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	Restore(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) Restore(id string) error {
	return r.db.Unscoped().Model(&models.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
