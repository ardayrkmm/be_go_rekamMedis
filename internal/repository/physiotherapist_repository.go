package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type PhysiotherapistRepository interface {
	FindAll(offset, limit int) ([]models.Physiotherapist, int64, error)
	FindByID(id string) (*models.Physiotherapist, error)
	Create(physio *models.Physiotherapist) error
	Update(physio *models.Physiotherapist) error
	Delete(id string) error
	Restore(id string) error
}

type physiotherapistRepository struct {
	db *firestore.Client
}

func NewPhysiotherapistRepository(db *firestore.Client) PhysiotherapistRepository {
	return &physiotherapistRepository{db}
}

func (r *physiotherapistRepository) FindAll(offset, limit int) ([]models.Physiotherapist, int64, error) {
	ctx := context.Background()
	var items []models.Physiotherapist
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("physiotherapists").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.Physiotherapist
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *physiotherapistRepository) FindByID(id string) (*models.Physiotherapist, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("physiotherapists").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.Physiotherapist
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *physiotherapistRepository) Create(item *models.Physiotherapist) error {
	ctx := context.Background()
	ref := r.db.Collection("physiotherapists").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *physiotherapistRepository) Update(item *models.Physiotherapist) error {
	ctx := context.Background()
	_, err := r.db.Collection("physiotherapists").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *physiotherapistRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("physiotherapists").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}

func (r *physiotherapistRepository) Restore(id string) error {
	ctx := context.Background()
	_, err := r.db.Collection("physiotherapists").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: nil},
	})
	return err
}
