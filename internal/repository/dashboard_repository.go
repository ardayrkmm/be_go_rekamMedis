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
		"pending":   0,
		"confirmed": 0,
		"completed": 0,
		"cancelled": 0,
	}
	for _, res := range results {
		statusMap[res.Status] = res.Count
	}

	var finalResult []map[string]interface{}
	for _, k := range []string{"pending", "confirmed", "completed", "cancelled"} {
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
