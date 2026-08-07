package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type PatientRepository interface {
	FindAll(offset, limit int) ([]models.Patient, int64, error)
	FindByID(id string) (*models.Patient, error)
	Create(patient *models.Patient) error
	Update(patient *models.Patient) error
	Delete(id string) error
	Restore(id string) error
}

type patientRepository struct {
	db *firestore.Client
}

func NewPatientRepository(db *firestore.Client) PatientRepository {
	return &patientRepository{db}
}

func (r *patientRepository) FindAll(offset, limit int) ([]models.Patient, int64, error) {
	ctx := context.Background()
	var patients []models.Patient
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("patients").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var patient models.Patient
		doc.DataTo(&patient)
		patient.ID = doc.Ref.ID
		patients = append(patients, patient)
	}

	return patients, total, nil
}

func (r *patientRepository) FindByID(id string) (*models.Patient, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("patients").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var patient models.Patient
	doc.DataTo(&patient)
	patient.ID = doc.Ref.ID
	return &patient, nil
}

func (r *patientRepository) Create(patient *models.Patient) error {
	ctx := context.Background()
	ref := r.db.Collection("patients").NewDoc()
	patient.ID = ref.ID
	_, err := ref.Set(ctx, patient)
	return err
}

func (r *patientRepository) Update(patient *models.Patient) error {
	ctx := context.Background()
	_, err := r.db.Collection("patients").Doc(patient.ID).Set(ctx, patient)
	return err
}

func (r *patientRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("patients").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}

func (r *patientRepository) Restore(id string) error {
	ctx := context.Background()
	_, err := r.db.Collection("patients").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: nil},
	})
	return err
}

