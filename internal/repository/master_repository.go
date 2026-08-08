package repository

import (
	"context"
	"backend_go/internal/models"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type MasterRepository interface {
	GetPatientCategories() ([]models.PatientCategory, error)
	GetGenders() ([]models.Gender, error)
	CreatePatientCategory(category *models.PatientCategory) error
}

type masterRepository struct {
	db *firestore.Client
}

func NewMasterRepository(db *firestore.Client) MasterRepository {
	return &masterRepository{db}
}

func (r *masterRepository) GetPatientCategories() ([]models.PatientCategory, error) {
	ctx := context.Background()
	var items []models.PatientCategory
	iter := r.db.Collection("patientcategories").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.PatientCategory
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *masterRepository) CreatePatientCategory(category *models.PatientCategory) error {
	ctx := context.Background()
	_, _, err := r.db.Collection("patientcategories").Add(ctx, category)
	return err
}

func (r *masterRepository) GetGenders() ([]models.Gender, error) {
	ctx := context.Background()
	var items []models.Gender
	iter := r.db.Collection("genders").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.Gender
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}
