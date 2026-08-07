package repository

import (
	"context"
	"backend_go/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type PaymentRepository interface {
	FindAll(offset, limit int) ([]models.Payment, int64, error)
	FindByID(id string) (*models.Payment, error)
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
}

type paymentRepository struct {
	db *firestore.Client
}

func NewPaymentRepository(db *firestore.Client) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) FindAll(offset, limit int) ([]models.Payment, int64, error) {
	ctx := context.Background()
	var items []models.Payment
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("payments").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.Payment
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *paymentRepository) FindByID(id string) (*models.Payment, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("payments").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.Payment
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *paymentRepository) Create(item *models.Payment) error {
	ctx := context.Background()
	ref := r.db.Collection("payments").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *paymentRepository) Update(item *models.Payment) error {
	ctx := context.Background()
	_, err := r.db.Collection("payments").Doc(item.ID).Set(ctx, item)
	return err
}
