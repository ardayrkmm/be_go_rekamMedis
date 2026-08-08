package usecase

import (
	"errors"
	"strings"
	"time"

	"backend_go/internal/models"
	"backend_go/internal/repository"
	"backend_go/pkg/utils"
)

type AuthUseCase interface {
	Login(email, password string) (string, *models.User, error)
	Register(user *models.User) error
	Profile(userID string) (*models.User, error)
	Logout(token string) error
	ForgotPassword(email string) error
	ResetPassword(email, token, newPassword string) error
	ChangePassword(userID string, oldPassword, newPassword string) error
}

type authUseCase struct {
	userRepo repository.UserRepository
	authRepo repository.AuthRepository
}

func NewAuthUseCase(userRepo repository.UserRepository, authRepo repository.AuthRepository) AuthUseCase {
	return &authUseCase{userRepo, authRepo}
}

func (u *authUseCase) Login(email, password string) (string, *models.User, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil || user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (u *authUseCase) Register(user *models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	if user.Role == "" {
		user.Role = string(models.RoleFisioterapis)
	}
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	return u.userRepo.Create(user)
}

func (u *authUseCase) Profile(userID string) (*models.User, error) {
	return u.userRepo.FindByID(userID)
}

func (u *authUseCase) Logout(token string) error {
	claims, err := utils.ValidateToken(token)
	if err != nil {
		return err
	}
	
	exp := int64(claims["exp"].(float64))
	expiresAt := time.Unix(exp, 0)

	blocklist := &models.JwtBlocklist{
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return u.authRepo.CreateBlocklist(blocklist)
}

func (u *authUseCase) ForgotPassword(email string) error {
	_, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	// Generate a random token, in real world we send it via email
	resetToken := &models.PasswordResetToken{
		Email:     email,
		Token:     "MOCK_RESET_TOKEN_123", // Mock token
		CreatedAt: time.Now(),
	}

	_ = u.authRepo.DeleteResetToken(email) // delete old if exists
	return u.authRepo.CreateResetToken(resetToken)
}

func (u *authUseCase) ResetPassword(email, token, newPassword string) error {
	_, err := u.authRepo.GetResetToken(email, token)
	if err != nil {
		return errors.New("invalid or expired token")
	}

	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := u.userRepo.Update(user); err != nil {
		return err
	}

	return u.authRepo.DeleteResetToken(email)
}

func (u *authUseCase) ChangePassword(userID string, oldPassword, newPassword string) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("incorrect old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return u.userRepo.Update(user)
}
