package repository

import (
	"context"
	"backend_go/internal/models"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type TherapySessionRepository interface {
	FindAll(offset, limit int) ([]models.TherapySession, int64, error)
	FindByID(id string) (*models.TherapySession, error)
	FindByAppointmentID(appointmentID string) ([]models.TherapySession, error)
	Create(session *models.TherapySession) error
	Update(session *models.TherapySession) error
	Delete(id string) error
	GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error)
	HasTreatedPatient(physioID, patientID string) (bool, error)
}

type therapySessionRepository struct {
	db *firestore.Client
}

func NewTherapySessionRepository(db *firestore.Client) TherapySessionRepository {
	return &therapySessionRepository{db}
}

func (r *therapySessionRepository) FindAll(offset, limit int) ([]models.TherapySession, int64, error) {
	ctx := context.Background()
	var items []models.TherapySession
	var total int64

	// total omitted for NoSQL

	iter := r.db.Collection("therapysessions").Where("DeletedAt", "==", nil).Offset(offset).Limit(limit).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, err
		}
		var item models.TherapySession
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		r.populateRelations(ctx, &item)
		items = append(items, item)
	}

	return items, total, nil
}

func (r *therapySessionRepository) FindByID(id string) (*models.TherapySession, error) {
	ctx := context.Background()
	doc, err := r.db.Collection("therapysessions").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var item models.TherapySession
	doc.DataTo(&item)
	item.ID = doc.Ref.ID
	r.populateRelations(ctx, &item)
	return &item, nil
}

func (r *therapySessionRepository) FindByAppointmentID(appointmentID string) ([]models.TherapySession, error) {
	ctx := context.Background()
	var items []models.TherapySession
	iter := r.db.Collection("therapysessions").Where("AppointmentID", "==", appointmentID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var item models.TherapySession
		doc.DataTo(&item)
		item.ID = doc.Ref.ID
		r.populateRelations(ctx, &item)
		items = append(items, item)
	}
	return items, nil
}

func (r *therapySessionRepository) Create(item *models.TherapySession) error {
	ctx := context.Background()
	ref := r.db.Collection("therapysessions").NewDoc()
	item.ID = ref.ID
	_, err := ref.Set(ctx, item)
	return err
}

func (r *therapySessionRepository) Update(item *models.TherapySession) error {
	ctx := context.Background()
	_, err := r.db.Collection("therapysessions").Doc(item.ID).Set(ctx, item)
	return err
}

func (r *therapySessionRepository) Delete(id string) error {
	ctx := context.Background()
	now := time.Now()
	_, err := r.db.Collection("therapysessions").Doc(id).Update(ctx, []firestore.Update{
		{Path: "DeletedAt", Value: &now},
	})
	return err
}

func (r *therapySessionRepository) GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error) {
	ctx := context.Background()
	var items []models.TherapySession
	
	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)

	iter := r.db.Collection("therapysessions").Where("TherapyDate", ">=", start).Where("TherapyDate", "<=", end).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done { break }
		if err != nil { return nil, err }
		var item models.TherapySession
		doc.DataTo(&item)
		if item.DeletedAt != nil { continue }
		item.ID = doc.Ref.ID
		r.populateRelations(ctx, &item)
		items = append(items, item)
	}
	return items, nil
}

func (r *therapySessionRepository) populateRelations(ctx context.Context, item *models.TherapySession) {
	if item.PatientID != "" {
		var pat models.Patient
		doc, err := r.db.Collection("patients").Doc(item.PatientID).Get(ctx)
		if err == nil {
			doc.DataTo(&pat)
			pat.ID = doc.Ref.ID
			item.Patient = &pat
		}
	}
	if item.PhysiotherapistID != "" {
		var phys models.Physiotherapist
		doc, err := r.db.Collection("physiotherapists").Doc(item.PhysiotherapistID).Get(ctx)
		if err == nil {
			doc.DataTo(&phys)
			phys.ID = doc.Ref.ID
			item.Physiotherapist = &phys
		}
	}
	if len(item.ServiceMasterIDs) > 0 {
		var srvs []models.ServiceMaster
		for _, id := range item.ServiceMasterIDs {
			var srv models.ServiceMaster
			doc, err := r.db.Collection("servicemasters").Doc(id).Get(ctx)
			if err == nil {
				doc.DataTo(&srv)
				srv.ID = doc.Ref.ID
				srvs = append(srvs, srv)
			}
		}
		item.ServiceMasters = srvs
		if len(srvs) > 0 {
			item.ServiceMaster = &srvs[0]
			if item.ServiceMasterID == "" {
				item.ServiceMasterID = srvs[0].ID
			}
		}
	} else if item.ServiceMasterID != "" {
		var srv models.ServiceMaster
		doc, err := r.db.Collection("servicemasters").Doc(item.ServiceMasterID).Get(ctx)
		if err == nil {
			doc.DataTo(&srv)
			srv.ID = doc.Ref.ID
			item.ServiceMaster = &srv
			item.ServiceMasters = []models.ServiceMaster{srv}
		}
	}
	if item.AppointmentID != "" {
		var apt models.Appointment
		doc, err := r.db.Collection("appointments").Doc(item.AppointmentID).Get(ctx)
		if err == nil {
			doc.DataTo(&apt)
			apt.ID = doc.Ref.ID
			item.Appointment = &apt
		}
	}
}

func (r *therapySessionRepository) HasTreatedPatient(physioID, patientID string) (bool, error) {
	ctx := context.Background()
	iter := r.db.Collection("therapysessions").Where("PhysiotherapistID", "==", physioID).Where("PatientID", "==", patientID).Where("DeletedAt", "==", nil).Limit(1).Documents(ctx)
	defer iter.Stop()
	
	_, err := iter.Next()
	if err == iterator.Done {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}


