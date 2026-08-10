package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type ActivityLogRepository interface {
	FindAll(offset, limit int) ([]models.ActivityLog, int64, error)
	FindByID(id string) (*models.ActivityLog, error)
	Delete(id string) error
}

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(db *gorm.DB) ActivityLogRepository {
	return &activityLogRepository{db}
}

func (r *activityLogRepository) FindAll(offset, limit int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	var total int64
	err := r.db.Model(&models.ActivityLog{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Preload("User").Offset(offset).Limit(limit).Find(&logs).Error
	return logs, total, err
}

func (r *activityLogRepository) FindByID(id string) (*models.ActivityLog, error) {
	var log models.ActivityLog
	err := r.db.Preload("User").First(&log, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *activityLogRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.ActivityLog{}).Error
}

