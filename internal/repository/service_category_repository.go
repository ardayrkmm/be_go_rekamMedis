package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ServiceCategoryRepository interface {
	FindAll(offset, limit int) ([]models.ServiceCategory, int64, error)
	FindByID(id string) (*models.ServiceCategory, error)
	Create(category *models.ServiceCategory) error
	Update(category *models.ServiceCategory) error
	Delete(id string) error
}

type serviceCategoryRepository struct {
	db *firestore.Client
}

func NewServiceCategoryRepository(db *firestore.Client) ServiceCategoryRepository {
	return &serviceCategoryRepository{db}
}

func (r *serviceCategoryRepository) FindAll(offset, limit int) ([]models.ServiceCategory, int64, error) {
	ctx := context.Background()
	var items []models.ServiceCategory
	var total int64

	iter := r.db.Collection("service_categories").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.ServiceCategory
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *serviceCategoryRepository) FindByID(id string) (*models.ServiceCategory, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("service_categories").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.ServiceCategory
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *serviceCategoryRepository) Create(item *models.ServiceCategory) error {
	ctx := context.Background()
	ref := r.db.Collection("service_categories").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *serviceCategoryRepository) Update(item *models.ServiceCategory) error {
	ctx := context.Background()
	_, err := r.db.Collection("service_categories").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *serviceCategoryRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("service_categories").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}
