package usecase

import (
	"backend_go/internal/repository"

	"github.com/gin-gonic/gin"
)

type DashboardUsecase interface {
	GetAdminDashboardData(filter string) (gin.H, error)
	GetFisioDashboardData(physiotherapistID string) (gin.H, error)
	GetReportsDashboardData() (gin.H, error)
}

type dashboardUsecase struct {
	repo repository.DashboardRepository
}

func NewDashboardUsecase(repo repository.DashboardRepository) DashboardUsecase {
	return &dashboardUsecase{repo: repo}
}

func (u *dashboardUsecase) GetAdminDashboardData(filter string) (gin.H, error) {
	// 1. Fetch counts
	totalPatients, _ := u.repo.CountPatients()
	totalFisio, _ := u.repo.CountPhysiotherapists()
	totalAppointments, _ := u.repo.CountAppointments()
	appointmentsToday, _ := u.repo.CountAppointmentsToday()
	newPatientsMonth, _ := u.repo.CountNewPatientsThisMonth()
	revenueToday, _ := u.repo.SumRevenueToday()
	revenueMonth, _ := u.repo.SumRevenueThisMonth()

	// 2. Build charts (For now, keep it dummy/mock since it requires complex aggregations that might not be needed immediately)
	charts := gin.H{
		"patients": []gin.H{
			{"label": "Jan", "total": 20},
			{"label": "Feb", "total": 35},
			{"label": "Mar", "total": totalPatients},
		},
		"appointments": []gin.H{
			{"label": "Jan", "total": 10},
			{"label": "Feb", "total": 20},
			{"label": "Mar", "total": totalAppointments},
		},
		"revenue": []gin.H{
			{"label": "Jan", "total": 5000000},
			{"label": "Feb", "total": 7500000},
			{"label": "Mar", "total": revenueMonth},
		},
		"diseases": []gin.H{
			{"diagnosis": "Low Back Pain", "count": 45},
			{"diagnosis": "Stroke", "count": 30},
			{"diagnosis": "Frozen Shoulder", "count": 25},
		},
	}

	return gin.H{
		"summary": gin.H{
			"total_pasien":          totalPatients,
			"total_fisioterapi":     totalFisio,
			"total_appointment":     totalAppointments,
			"appointment_hari_ini":  appointmentsToday,
			"pasien_baru_bulan_ini": newPatientsMonth,
			"pendapatan_hari_ini":   revenueToday,
			"pendapatan_bulan_ini":  revenueMonth,
		},
		"charts": charts,
	}, nil
}

func (u *dashboardUsecase) GetFisioDashboardData(physiotherapistID string) (gin.H, error) {
	// Simplistic approach, for fisio dashboard we can fetch counts
	totalPatients, _ := u.repo.CountPatients() // Ideally should be patients for this physio
	appointmentsToday, _ := u.repo.CountAppointmentsToday()
	
	return gin.H{
		"physiotherapist_name": "Fisioterapis",
		"today_appointments": appointmentsToday,
		"today_patients": totalPatients,
		"today_therapy_sessions": 0,
		"next_appointment": gin.H{
			"time": "14:30",
			"patient_name": "Pasien",
		},
		"recent_activities": []gin.H{},
	}, nil
}

func (u *dashboardUsecase) GetReportsDashboardData() (gin.H, error) {
	totalPatients, _ := u.repo.CountPatients()
	totalFisio, _ := u.repo.CountPhysiotherapists()
	totalAppointments, _ := u.repo.CountAppointments()
	
	// Count total sessions (dummy for now as we don't have CountSessions)
	totalSessions := 120 

	// Dummy data for appointments by status
	appointmentsByStatus := []gin.H{
		{"status": "pending", "count": 15},
		{"status": "confirmed", "count": 45},
		{"status": "completed", "count": 30},
		{"status": "cancelled", "count": 10},
	}

	// Dummy data for sessions by month
	sessionsByMonth := []gin.H{
		{"month": "2026-02", "count": 10},
		{"month": "2026-03", "count": 25},
		{"month": "2026-04", "count": 40},
		{"month": "2026-05", "count": 60},
		{"month": "2026-06", "count": 55},
		{"month": "2026-07", "count": 70},
	}

	return gin.H{
		"summary": gin.H{
			"total_patients": totalPatients,
			"total_physiotherapists": totalFisio,
			"total_appointments": totalAppointments,
			"total_sessions": totalSessions,
		},
		"appointments_by_status": appointmentsByStatus,
		"sessions_by_month": sessionsByMonth,
	}, nil
}
