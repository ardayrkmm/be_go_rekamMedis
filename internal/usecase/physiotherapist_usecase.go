package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"backend_go/pkg/utils"
	"strings"
	"time"
)

type PhysiotherapistUseCase interface {
	Fetch(offset, limit int, search string) ([]models.Physiotherapist, int64, error)
	GetByID(id string) (*models.Physiotherapist, error)
	Store(physio *models.Physiotherapist) error
	Update(id string, req *models.Physiotherapist) error
	Delete(id string) error
	Restore(id string) error
}

type physiotherapistUseCase struct {
	physioRepo repository.PhysiotherapistRepository
	userRepo   repository.UserRepository
}

func NewPhysiotherapistUseCase(physioRepo repository.PhysiotherapistRepository, userRepo repository.UserRepository) PhysiotherapistUseCase {
	return &physiotherapistUseCase{
		physioRepo: physioRepo,
		userRepo:   userRepo,
	}
}

func (u *physiotherapistUseCase) Fetch(offset, limit int, search string) ([]models.Physiotherapist, int64, error) {
	return u.physioRepo.FindAll(offset, limit, search)
}

func (u *physiotherapistUseCase) GetByID(id string) (*models.Physiotherapist, error) {
	return u.physioRepo.FindByID(id)
}

func (u *physiotherapistUseCase) Store(physio *models.Physiotherapist) error {
	physio.CreatedAt = time.Now()
	physio.UpdatedAt = time.Now()
	if physio.Status == "" {
		physio.Status = "Aktif"
	}
	if physio.Sip != nil && *physio.Sip == "" {
		physio.Sip = nil
	}

	if physio.Password != "" && physio.Email != "" {
		hashedPassword, err := utils.HashPassword(physio.Password)
		if err != nil {
			return err
		}
		user := &models.User{
			Name:     physio.Name,
			Email:    strings.ToLower(strings.TrimSpace(physio.Email)),
			Password: hashedPassword,
			Role:     string(models.RoleFisioterapis),
		}
		if err := u.userRepo.Create(user); err != nil {
			return err
		}
	}

	return u.physioRepo.Create(physio)
}

func (u *physiotherapistUseCase) Update(id string, req *models.Physiotherapist) error {
	physio, err := u.physioRepo.FindByID(id)
	if err != nil {
		return err
	}

	if physio.Email != "" {
		emailToFind := strings.ToLower(strings.TrimSpace(physio.Email))
		user, errUser := u.userRepo.FindByEmail(emailToFind)
		if errUser != nil || user == nil {
			user, errUser = u.userRepo.FindByEmail(physio.Email)
		}
		if errUser == nil && user != nil && user.Role == string(models.RoleFisioterapis) {
			user.Name = req.Name
			user.Email = strings.ToLower(strings.TrimSpace(req.Email))
			_ = u.userRepo.Update(user)
		}
	}

	physio.Name = req.Name
	physio.Specialization = req.Specialization
	if req.Sip != nil && *req.Sip == "" {
		physio.Sip = nil
	} else {
		physio.Sip = req.Sip
	}
	physio.Phone = req.Phone
	physio.Email = req.Email
	physio.Address = req.Address
	physio.Gender = req.Gender
	if req.Photo != nil {
		physio.Photo = req.Photo
	}
	physio.Status = req.Status
	physio.UpdatedAt = time.Now()

	return u.physioRepo.Update(physio)
}

func (u *physiotherapistUseCase) Delete(id string) error {
	physio, err := u.physioRepo.FindByID(id)
	if err == nil && physio.Email != "" {
		emailToFind := strings.ToLower(strings.TrimSpace(physio.Email))
		user, errUser := u.userRepo.FindByEmail(emailToFind)
		if errUser != nil || user == nil {
			user, errUser = u.userRepo.FindByEmail(physio.Email)
		}
		if errUser == nil && user != nil && user.Role == string(models.RoleFisioterapis) {
			_ = u.userRepo.Delete(user.ID)
		}
	}
	return u.physioRepo.Delete(id)
}

func (u *physiotherapistUseCase) Restore(id string) error {
	return u.physioRepo.Restore(id)
}
