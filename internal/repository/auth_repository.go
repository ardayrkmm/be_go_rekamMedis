package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type AuthRepository interface {
	CreateResetToken(token *models.PasswordResetToken) error
	GetResetToken(email, token string) (*models.PasswordResetToken, error)
	DeleteResetToken(email string) error
	CreateBlocklist(blocklist *models.JwtBlocklist) error
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db}
}

func (r *authRepository) CreateResetToken(token *models.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *authRepository) GetResetToken(email, token string) (*models.PasswordResetToken, error) {
	var resetToken models.PasswordResetToken
	err := r.db.Where("email = ? AND token = ?", email, token).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *authRepository) DeleteResetToken(email string) error {
	return r.db.Where("email = ?", email).Delete(&models.PasswordResetToken{}).Error
}

func (r *authRepository) CreateBlocklist(blocklist *models.JwtBlocklist) error {
	return r.db.Create(blocklist).Error
}
