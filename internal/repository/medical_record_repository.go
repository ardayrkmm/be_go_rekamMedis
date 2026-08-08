package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type MedicalRecordRepository interface {
	FindAll(offset, limit int) ([]models.MedicalRecord, int64, error)
	FindByID(id string) (*models.MedicalRecord, error)
	FindByPatientID(patientID string) ([]models.MedicalRecord, error)
	Create(record *models.MedicalRecord) error
	Update(record *models.MedicalRecord) error
	Delete(id string) error
}

type medicalRecordRepository struct {
	db *firestore.Client
}

func NewMedicalRecordRepository(db *firestore.Client) MedicalRecordRepository {
	return &medicalRecordRepository{db}
}

func (r *medicalRecordRepository) FindAll(offset, limit int) ([]models.MedicalRecord, int64, error) {
	ctx := context.Background()
	var items []models.MedicalRecord
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("medicalrecords").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.MedicalRecord
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}

	return items, total, nil
}

func (r *medicalRecordRepository) FindByID(id string) (*models.MedicalRecord, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("medicalrecords").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.MedicalRecord
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *medicalRecordRepository) FindByPatientID(patientID string) ([]models.MedicalRecord, error) {
	ctx := context.Background()
	var items []models.MedicalRecord
	iter := r.db.Collection("medicalrecords").Where("PatientID", "==", patientID).Where("DeletedAt", "==", nil).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item models.MedicalRecord
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *medicalRecordRepository) Create(item *models.MedicalRecord) error {
	ctx := context.Background()
	ref := r.db.Collection("medicalrecords").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *medicalRecordRepository) Update(item *models.MedicalRecord) error {
	ctx := context.Background()
	_, err := r.db.Collection("medicalrecords").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *medicalRecordRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("medicalrecords").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}
