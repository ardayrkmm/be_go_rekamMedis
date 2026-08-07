package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type AppointmentRepository interface {
	FindAll(offset, limit int) ([]models.Appointment, int64, error)
	FindByID(id string) (*models.Appointment, error)
	FindByPatientID(patientID string) ([]models.Appointment, error)
	FindByPhysiotherapistID(physioID string) ([]models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(appointment *models.Appointment) error
	Delete(id string) error
}

type appointmentRepository struct {
	db *firestore.Client
}

func NewAppointmentRepository(db *firestore.Client) AppointmentRepository {
	return &appointmentRepository{db}
}

func (r *appointmentRepository) FindAll(offset, limit int) ([]models.Appointment, int64, error) {
	ctx := context.Background()
	var items []models.Appointment
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("appointments").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.Appointment
		doc.DataTo(&item)
		item.ID = doc.Ref.ID

		if item.PatientID != "" {
			var pat models.Patient
			patDoc, err := r.db.Collection("patients").Doc(item.PatientID).Get(ctx)
			if err == nil {
				patDoc.DataTo(&pat)
				pat.ID = patDoc.Ref.ID
				item.Patient = &pat
			}
		}

		if item.PhysiotherapistID != "" {
			var physio models.Physiotherapist
			phyDoc, err := r.db.Collection("physiotherapists").Doc(item.PhysiotherapistID).Get(ctx)
			if err == nil {
				phyDoc.DataTo(&physio)
				physio.ID = phyDoc.Ref.ID
				item.Physiotherapist = &physio
			}
		}

		if item.ServiceMasterID != "" {
			var sm models.ServiceMaster
			smDoc, err := r.db.Collection("servicemasters").Doc(item.ServiceMasterID).Get(ctx)
			if err == nil {
				smDoc.DataTo(&sm)
				sm.ID = smDoc.Ref.ID
				item.ServiceMaster = &sm
			}
		}
		items = append(items, item)
	}

	return items, total, nil
}

func (r *appointmentRepository) FindByID(id string) (*models.Appointment, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("appointments").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.Appointment
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	return &item, nil
}

func (r *appointmentRepository) FindByPatientID(patientID string) ([]models.Appointment, error) {
	ctx := context.Background()
	var items []models.Appointment
	iter := r.db.Collection("appointments").Where("PatientID", "==", patientID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item models.Appointment
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *appointmentRepository) FindByPhysiotherapistID(physiotherapistID string) ([]models.Appointment, error) {
	ctx := context.Background()
	var items []models.Appointment
	iter := r.db.Collection("appointments").Where("PhysiotherapistID", "==", physiotherapistID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item models.Appointment
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		items = append(items, item)
	}
	return items, nil
}

func (r *appointmentRepository) Create(item *models.Appointment) error {
	ctx := context.Background()
	ref := r.db.Collection("appointments").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *appointmentRepository) Update(item *models.Appointment) error {
	ctx := context.Background()
	_, err := r.db.Collection("appointments").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *appointmentRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("appointments").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}
