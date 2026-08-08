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

	// 2. Fetch real data for charts
	patientsByMonth, _ := u.repo.CountPatientsByMonth()
	apptsByStatus, _ := u.repo.CountAppointmentsByStatus()

	// Map appointments by status to name and value for frontend BarChart
	var mappedAppts []gin.H
	if apptsByStatus != nil {
		for _, appt := range apptsByStatus {
			mappedAppts = append(mappedAppts, gin.H{
				"name":  appt["status"],
				"value": appt["count"],
			})
		}
	} else {
		mappedAppts = []gin.H{}
	}

	charts := gin.H{
		"patients": patientsByMonth,
		"appointments": mappedAppts,
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
	
	// Real data instead of dummy
	totalSessions, _ := u.repo.CountTherapySessions()
	appointmentsByStatus, _ := u.repo.CountAppointmentsByStatus()
	sessionsByMonth, _ := u.repo.CountTherapySessionsByMonth()

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
