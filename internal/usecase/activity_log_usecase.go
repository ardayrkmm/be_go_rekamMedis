package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
)

type ActivityLogUseCase interface {
	Fetch(offset, limit int) ([]models.ActivityLog, int64, error)
	GetByID(id string) (*models.ActivityLog, error)
	Delete(id string) error
}

type activityLogUseCase struct {
	logRepo repository.ActivityLogRepository
}

func NewActivityLogUseCase(logRepo repository.ActivityLogRepository) ActivityLogUseCase {
	return &activityLogUseCase{
		logRepo: logRepo,
	}
}

func (u *activityLogUseCase) Fetch(offset, limit int) ([]models.ActivityLog, int64, error) {
	return u.logRepo.FindAll(offset, limit)
}

func (u *activityLogUseCase) GetByID(id string) (*models.ActivityLog, error) {
	return u.logRepo.FindByID(id)
}

func (u *activityLogUseCase) Delete(id string) error {
	return u.logRepo.Delete(id)
}
