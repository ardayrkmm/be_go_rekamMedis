package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ActivityLogRepository interface {
	FindAll(offset, limit int) ([]models.ActivityLog, int64, error)
	FindByID(id string) (*models.ActivityLog, error)
	Delete(id string) error
}

type activityLogRepository struct {
	db *firestore.Client
}

func NewActivityLogRepository(db *firestore.Client) ActivityLogRepository {
	return &activityLogRepository{db}
}

func (r *activityLogRepository) FindAll(offset, limit int) ([]models.ActivityLog, int64, error) {
	ctx := context.Background()
	var items []models.ActivityLog
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("activitylogs").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.ActivityLog
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *activityLogRepository) FindByID(id string) (*models.ActivityLog, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("activitylogs").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.ActivityLog
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *activityLogRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("activitylogs").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}
