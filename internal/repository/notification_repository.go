package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
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
	db *firestore.Client
}

func NewNotificationRepository(db *firestore.Client) NotificationRepository {
	return &notificationRepository{db}
}

func (r *notificationRepository) FindAll(offset, limit int) ([]models.Notification, int64, error) {
	ctx := context.Background()
	var items []models.Notification
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("notifications").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.Notification
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *notificationRepository) FindByID(id string) (*models.Notification, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("notifications").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.Notification
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *notificationRepository) Create(item *models.Notification) error {
	ctx := context.Background()
	ref := r.db.Collection("notifications").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *notificationRepository) Update(item *models.Notification) error {
	ctx := context.Background()
	_, err := r.db.Collection("notifications").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *notificationRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("notifications").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}

func (r *notificationRepository) FindUnread() ([]models.Notification, error) {
	ctx := context.Background()
	var items []models.Notification
	iter := r.db.Collection("notifications").Where("IsRead", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.Notification
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *notificationRepository) MarkAllAsRead() error {
	ctx := context.Background()
	iter := r.db.Collection("notifications").Where("IsRead", "==", false).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return err }
		doc.Ref.Update(ctx, []firestore.Update{{Path: "IsRead", Value: true}})
	}
	return nil
}
