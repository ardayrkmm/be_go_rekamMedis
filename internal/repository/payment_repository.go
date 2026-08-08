package repository

import (
	"context"
	"strings"
	"sort"
	"time"
	"backend_go/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type PaymentRepository interface {
	FindAll(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error)
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

func (r *paymentRepository) FindAll(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error) {
	ctx := context.Background()
	var allItems []models.Payment

	// Fetch all undeleted documents
	iter := r.db.Collection("payments").Where("DeletedAt", "==", nil).Documents(ctx)
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
		allItems = append(allItems, item)
	}

	// Filter in Go
	var filtered []models.Payment
	var start, end time.Time
	hasDateFilter := false
	if startDate != "" && endDate != "" {
		if s, err := time.Parse("2006-01-02", startDate); err == nil {
			if e, err := time.Parse("2006-01-02", endDate); err == nil {
				start = s
				end = e.Add(24 * time.Hour).Add(-time.Nanosecond)
				hasDateFilter = true
			}
		}
	}

	for _, item := range allItems {
		// Filter by status
		if status != "" && status != "all" && item.Status != status {
			continue
		}
		// Filter by date
		if hasDateFilter && item.PaymentDate != nil {
			if item.PaymentDate.Before(start) || item.PaymentDate.After(end) {
				continue
			}
		}
		// Filter by search (case-insensitive simple match on invoice/patient)
		if search != "" {
			searchLower := strings.ToLower(search)
			if !strings.Contains(strings.ToLower(item.InvoiceNumber), searchLower) &&
				!strings.Contains(strings.ToLower(item.PatientName), searchLower) &&
				!strings.Contains(strings.ToLower(item.PhysiotherapistName), searchLower) {
				continue
			}
		}
		filtered = append(filtered, item)
	}

	// Sort by created at descending (latest first)
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
	})

	total := int64(len(filtered))

	// Paginate
	if offset >= len(filtered) {
		return []models.Payment{}, total, nil
	}
	endIdx := offset + limit
	if endIdx > len(filtered) {
		endIdx = len(filtered)
	}

	return filtered[offset:endIdx], total, nil
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
