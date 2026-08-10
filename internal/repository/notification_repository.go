package repository

import (
	"time"
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	FindAll(offset, limit int) ([]models.Notification, int64, error)
	FindUnread() ([]models.Notification, error)
	FindByID(id string) (*models.Notification, error)
	Create(notification *models.Notification) error
	Update(notification *models.Notification) error
	Delete(id string) error
	MarkAllAsRead() error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) FindAll(offset, limit int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64
	err := r.db.Model(&models.Notification{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.Offset(offset).Limit(limit).Find(&notifications).Error
	return notifications, total, err
}

func (r *notificationRepository) FindUnread() ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("read_at IS NULL").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindByID(id string) (*models.Notification, error) {
	var notification models.Notification
	err := r.db.First(&notification, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &notification, nil
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) Update(notification *models.Notification) error {
	return r.db.Save(notification).Error
}

func (r *notificationRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Notification{}).Error
}

func (r *notificationRepository) MarkAllAsRead() error {
	now := time.Now()
	return r.db.Model(&models.Notification{}).Where("read_at IS NULL").Update("read_at", now).Error
}

