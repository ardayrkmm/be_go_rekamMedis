package repository

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type DashboardRepository interface {
	CountPatients() (int, error)
	CountPhysiotherapists() (int, error)
	CountAppointments() (int, error)
	CountAppointmentsToday() (int, error)
	CountNewPatientsThisMonth() (int, error)
	SumRevenueToday() (int, error)
	SumRevenueThisMonth() (int, error)
}

type dashboardRepository struct {
	db *firestore.Client
}

func NewDashboardRepository(db *firestore.Client) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) countDocs(iter *firestore.DocumentIterator) (int, error) {
	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil
}

func (r *dashboardRepository) CountPatients() (int, error) {
	ctx := context.Background()
	iter := r.db.Collection("patients").Where("DeletedAt", "==", nil).Documents(ctx)
	return r.countDocs(iter)
}

func (r *dashboardRepository) CountPhysiotherapists() (int, error) {
	ctx := context.Background()
	iter := r.db.Collection("users").Where("Role", "==", "physiotherapist").Where("DeletedAt", "==", nil).Documents(ctx)
	return r.countDocs(iter)
}

func (r *dashboardRepository) CountAppointments() (int, error) {
	ctx := context.Background()
	iter := r.db.Collection("appointments").Where("DeletedAt", "==", nil).Documents(ctx)
	return r.countDocs(iter)
}

func (r *dashboardRepository) CountAppointmentsToday() (int, error) {
	ctx := context.Background()
	todayStr := time.Now().Format("2006-01-02")
	iter := r.db.Collection("appointments").Where("AppointmentDate", "==", todayStr).Where("DeletedAt", "==", nil).Documents(ctx)
	return r.countDocs(iter)
}

func (r *dashboardRepository) CountNewPatientsThisMonth() (int, error) {
	// For simplicity, just return 0 or do a simple query if CreatedAt exists
	// We will just return 0 for now as it's complex without proper dates in dummy data
	return 0, nil
}

func (r *dashboardRepository) SumRevenueToday() (int, error) {
	ctx := context.Background()
	todayStr := time.Now().Format("2006-01-02")
	
	// We need to parse PaymentDate string since we store it as string or timestamp
	iter := r.db.Collection("payments").Where("Status", "==", "paid").Documents(ctx)
	
	sum := 0.0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		
		data := doc.Data()
		
		// Simplistic filtering by prefix of string if it's ISO8601
		if dateStr, ok := data["PaymentDate"].(string); ok {
			if len(dateStr) >= 10 && dateStr[:10] == todayStr {
				if total, ok := data["Total"].(float64); ok {
					sum += total
				}
			}
		} else if dateTime, ok := data["PaymentDate"].(time.Time); ok {
			if dateTime.Format("2006-01-02") == todayStr {
				if total, ok := data["Total"].(float64); ok {
					sum += total
				}
			}
		}
	}
	return int(sum), nil
}

func (r *dashboardRepository) SumRevenueThisMonth() (int, error) {
	ctx := context.Background()
	monthStr := time.Now().Format("2006-01")
	
	iter := r.db.Collection("payments").Where("Status", "==", "paid").Documents(ctx)
	
	sum := 0.0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		
		data := doc.Data()
		
		if dateStr, ok := data["PaymentDate"].(string); ok {
			if len(dateStr) >= 7 && dateStr[:7] == monthStr {
				if total, ok := data["Total"].(float64); ok {
					sum += total
				}
			}
		} else if dateTime, ok := data["PaymentDate"].(time.Time); ok {
			if dateTime.Format("2006-01") == monthStr {
				if total, ok := data["Total"].(float64); ok {
					sum += total
				}
			}
		}
	}
	return int(sum), nil
}
