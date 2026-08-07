package usecase

import (
	"time"
	"backend_go/internal/models"
	"backend_go/internal/repository"
)

type NotificationUseCase interface {
	Fetch(offset, limit int) ([]models.Notification, int64, error)
	FetchUnread() ([]models.Notification, error)
	MarkAllAsRead() error
	MarkAsRead(id string) error
	MarkAsUnread(id string) error
	Delete(id string) error
}

type notificationUseCase struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationUseCase(notificationRepo repository.NotificationRepository) NotificationUseCase {
	return &notificationUseCase{
		notificationRepo: notificationRepo,
	}
}

func (u *notificationUseCase) Fetch(offset, limit int) ([]models.Notification, int64, error) {
	return u.notificationRepo.FindAll(offset, limit)
}

func (u *notificationUseCase) FetchUnread() ([]models.Notification, error) {
	return u.notificationRepo.FindUnread()
}

func (u *notificationUseCase) MarkAllAsRead() error {
	return u.notificationRepo.MarkAllAsRead()
}

func (u *notificationUseCase) MarkAsRead(id string) error {
	notification, err := u.notificationRepo.FindByID(id)
	if err != nil {
		return err
	}
	now := time.Now()
	notification.ReadAt = &now
	return u.notificationRepo.Update(notification)
}

func (u *notificationUseCase) MarkAsUnread(id string) error {
	notification, err := u.notificationRepo.FindByID(id)
	if err != nil {
		return err
	}
	notification.ReadAt = nil
	return u.notificationRepo.Update(notification)
}

func (u *notificationUseCase) Delete(id string) error {
	return u.notificationRepo.Delete(id)
}
