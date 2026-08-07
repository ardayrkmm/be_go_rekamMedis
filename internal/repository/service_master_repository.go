package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type ServiceMasterRepository interface {
	FindAll(offset, limit int) ([]models.ServiceMaster, int64, error)
	FindByID(id string) (*models.ServiceMaster, error)
	Create(service *models.ServiceMaster) error
	Update(service *models.ServiceMaster) error
	Delete(id string) error
}

type serviceMasterRepository struct {
	db *firestore.Client
}

func NewServiceMasterRepository(db *firestore.Client) ServiceMasterRepository {
	return &serviceMasterRepository{db}
}

func (r *serviceMasterRepository) FindAll(offset, limit int) ([]models.ServiceMaster, int64, error) {
	ctx := context.Background()
	var items []models.ServiceMaster
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("servicemasters").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.ServiceMaster
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *serviceMasterRepository) FindByID(id string) (*models.ServiceMaster, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("servicemasters").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.ServiceMaster
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *serviceMasterRepository) Create(item *models.ServiceMaster) error {
	ctx := context.Background()
	ref := r.db.Collection("servicemasters").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *serviceMasterRepository) Update(item *models.ServiceMaster) error {
	ctx := context.Background()
	_, err := r.db.Collection("servicemasters").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *serviceMasterRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("servicemasters").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}
