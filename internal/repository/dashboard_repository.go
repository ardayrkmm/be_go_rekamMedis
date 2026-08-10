package repository

import (
	"time"

	"backend_go/internal/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	CountPatients() (int, error)
	CountPhysiotherapists() (int, error)
	CountAppointments() (int, error)
	CountAppointmentsToday() (int, error)
	CountNewPatientsThisMonth() (int, error)
	SumRevenueToday() (int, error)
	SumRevenueThisMonth() (int, error)
	CountTherapySessions() (int, error)
	CountAppointmentsByStatus() ([]map[string]interface{}, error)
	CountTherapySessionsByMonth() ([]map[string]interface{}, error)
	CountPatientsByMonth() ([]map[string]interface{}, error)

	GetPhysioIDByUserID(userID string) (string, error)
	CountPhysioPatients(physioID string) (int, error)
	CountPhysioAppointmentsToday(physioID string) (int, error)
	CountPhysioTherapySessionsToday(physioID string) (int, error)
	CountPhysioPatientsByMonth(physioID string) ([]map[string]interface{}, error)
	CountPhysioAppointmentsByStatus(physioID string) ([]map[string]interface{}, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

func (r *dashboardRepository) CountPatients() (int, error) {
	var count int64
	err := r.db.Model(&models.Patient{}).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountPhysiotherapists() (int, error) {
	var count int64
	err := r.db.Model(&models.Physiotherapist{}).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountAppointments() (int, error) {
	var count int64
	err := r.db.Model(&models.Appointment{}).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountAppointmentsToday() (int, error) {
	var count int64
	// In GORM, if AppointmentDate is time.Time, it's easier to query by range.
	// We'll just do a basic string query if it's a string, or range if time.
	// Since AppointmentDate is time.Time, we should compare dates.
	todayStr := time.Now().Format("2006-01-02")
	err := r.db.Model(&models.Appointment{}).Where("DATE(appointment_date) = ?", todayStr).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountNewPatientsThisMonth() (int, error) {
	return 0, nil
}

func (r *dashboardRepository) SumRevenueToday() (int, error) {
	var total float64
	todayStr := time.Now().Format("2006-01-02")
	err := r.db.Model(&models.Payment{}).Where("status = ? AND DATE(payment_date) = ?", "paid", todayStr).Select("IFNULL(SUM(total), 0)").Scan(&total).Error
	return int(total), err
}

func (r *dashboardRepository) SumRevenueThisMonth() (int, error) {
	var total float64
	monthStr := time.Now().Format("2006-01")
	err := r.db.Model(&models.Payment{}).Where("status = ? AND DATE_FORMAT(payment_date, '%Y-%m') = ?", "paid", monthStr).Select("IFNULL(SUM(total), 0)").Scan(&total).Error
	return int(total), err
}

func (r *dashboardRepository) CountTherapySessions() (int, error) {
	var count int64
	err := r.db.Model(&models.TherapySession{}).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountAppointmentsByStatus() ([]map[string]interface{}, error) {
	type Result struct {
		Status string
		Count  int
	}
	var results []Result
	err := r.db.Model(&models.Appointment{}).Select("status, COUNT(id) as count").Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	statusMap := map[string]int{
		"scheduled":   0,
		"telah_tiba": 0,
		"ongoing": 0,
		"completed": 0,
		"cancelled": 0,
	}
	for _, res := range results {
		statusMap[res.Status] = res.Count
	}

	var finalResult []map[string]interface{}
	for _, k := range []string{"scheduled", "telah_tiba", "ongoing", "completed", "cancelled"} {
		finalResult = append(finalResult, map[string]interface{}{
			"status": k,
			"count":  statusMap[k],
		})
	}
	return finalResult, nil
}

func (r *dashboardRepository) CountTherapySessionsByMonth() ([]map[string]interface{}, error) {
	type Result struct {
		Month string
		Count int
	}
	var results []Result
	sixMonthsAgo := time.Now().AddDate(0, -5, 0)
	startOfMonth := time.Date(sixMonthsAgo.Year(), sixMonthsAgo.Month(), 1, 0, 0, 0, 0, sixMonthsAgo.Location())

	err := r.db.Model(&models.TherapySession{}).
		Select("DATE_FORMAT(therapy_date, '%Y-%m') as month, COUNT(id) as count").
		Where("therapy_date >= ?", startOfMonth).
		Group("month").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	monthCounts := make(map[string]int)
	for _, res := range results {
		monthCounts[res.Month] = res.Count
	}

	var finalResult []map[string]interface{}
	for i := 5; i >= 0; i-- {
		m := time.Now().AddDate(0, -i, 0).Format("2006-01")
		finalResult = append(finalResult, map[string]interface{}{
			"month": m,
			"count": monthCounts[m],
		})
	}

	return finalResult, nil
}

func (r *dashboardRepository) CountPatientsByMonth() ([]map[string]interface{}, error) {
	type Result struct {
		Month string
		Count int
	}
	var results []Result
	sixMonthsAgo := time.Now().AddDate(0, -5, 0)
	startOfMonth := time.Date(sixMonthsAgo.Year(), sixMonthsAgo.Month(), 1, 0, 0, 0, 0, sixMonthsAgo.Location())

	err := r.db.Model(&models.Patient{}).
		Select("DATE_FORMAT(created_at, '%Y-%m') as month, COUNT(id) as count").
		Where("created_at >= ?", startOfMonth).
		Group("month").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	monthCounts := make(map[string]int)
	monthsMap := map[string]string{
		"01": "Jan", "02": "Feb", "03": "Mar", "04": "Apr", "05": "May", "06": "Jun",
		"07": "Jul", "08": "Aug", "09": "Sep", "10": "Oct", "11": "Nov", "12": "Dec",
	}

	for _, res := range results {
		monthCounts[res.Month] = res.Count
	}

	var finalResult []map[string]interface{}
	for i := 5; i >= 0; i-- {
		targetTime := time.Now().AddDate(0, -i, 0)
		mStr := targetTime.Format("2006-01")
		monthNum := targetTime.Format("01")
		finalResult = append(finalResult, map[string]interface{}{
			"month": monthsMap[monthNum],
			"total": monthCounts[mStr],
		})
	}

	return finalResult, nil
}

func (r *dashboardRepository) GetPhysioIDByUserID(userID string) (string, error) {
	// Step 1: get the user's email from the users table
	var userEmail string
	err := r.db.Table("users").Where("id = ?", userID).Select("email").Scan(&userEmail).Error
	if err != nil || userEmail == "" {
		return "", err
	}

	// Step 2: find physiotherapist by email
	var physioID string
	err = r.db.Table("physiotherapists").Where("email = ?", userEmail).Select("id").Scan(&physioID).Error
	return physioID, err
}

func (r *dashboardRepository) CountPhysioPatients(physioID string) (int, error) {
	var count int64
	err := r.db.Model(&models.Appointment{}).Where("physiotherapist_id = ?", physioID).Select("COUNT(DISTINCT patient_id)").Scan(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountPhysioAppointmentsToday(physioID string) (int, error) {
	var count int64
	todayStr := time.Now().Format("2006-01-02")
	err := r.db.Model(&models.Appointment{}).Where("physiotherapist_id = ? AND DATE(appointment_date) = ?", physioID, todayStr).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountPhysioTherapySessionsToday(physioID string) (int, error) {
	var count int64
	todayStr := time.Now().Format("2006-01-02")
	err := r.db.Model(&models.TherapySession{}).Joins("JOIN appointments on appointments.id = therapy_sessions.appointment_id").Where("appointments.physiotherapist_id = ? AND DATE(therapy_sessions.therapy_date) = ?", physioID, todayStr).Count(&count).Error
	return int(count), err
}

func (r *dashboardRepository) CountPhysioPatientsByMonth(physioID string) ([]map[string]interface{}, error) {
	type Result struct {
		Month string
		Count int
	}
	var results []Result
	sixMonthsAgo := time.Now().AddDate(0, -5, 0)
	startOfMonth := time.Date(sixMonthsAgo.Year(), sixMonthsAgo.Month(), 1, 0, 0, 0, 0, sixMonthsAgo.Location())

	// For patients by month for this physio, we count unique patients from appointments
	err := r.db.Model(&models.Appointment{}).
		Select("DATE_FORMAT(appointment_date, '%Y-%m') as month, COUNT(DISTINCT patient_id) as count").
		Where("physiotherapist_id = ? AND appointment_date >= ?", physioID, startOfMonth).
		Group("month").Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Prepare last 6 months list
	var finalResults []map[string]interface{}
	for i := 5; i >= 0; i-- {
		m := time.Now().AddDate(0, -i, 0)
		monthStr := m.Format("2006-01")
		monthName := m.Format("Jan")

		count := 0
		for _, r := range results {
			if r.Month == monthStr {
				count = r.Count
				break
			}
		}

		finalResults = append(finalResults, map[string]interface{}{
			"month": monthName,
			"total": count,
		})
	}
	return finalResults, nil
}

func (r *dashboardRepository) CountPhysioAppointmentsByStatus(physioID string) ([]map[string]interface{}, error) {
	type Result struct {
		Status string
		Count  int
	}
	var results []Result
	err := r.db.Model(&models.Appointment{}).Where("physiotherapist_id = ?", physioID).Select("status, COUNT(id) as count").Group("status").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	statusMap := map[string]int{
		"scheduled":  0,
		"telah_tiba": 0,
		"ongoing":    0,
		"completed":  0,
		"cancelled":  0,
	}

	for _, r := range results {
		statusMap[r.Status] = r.Count
	}

	return []map[string]interface{}{
		{"name": "Selesai", "value": statusMap["completed"]},
		{"name": "Dijadwalkan", "value": statusMap["scheduled"] + statusMap["telah_tiba"] + statusMap["ongoing"]},
		{"name": "Dibatalkan", "value": statusMap["cancelled"]},
	}, nil
}
