package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type AppointmentRepository interface {
	FindAll(offset, limit int, search string, status string, startDate string, endDate string) ([]models.Appointment, int64, error)
	FindByID(id string) (*models.Appointment, error)
	FindByPatientID(patientID string) ([]models.Appointment, error)
	FindByPhysiotherapistID(physioID string) ([]models.Appointment, error)
	Create(appointment *models.Appointment) error
	Update(appointment *models.Appointment) error
	Delete(id string) error
}

type appointmentRepository struct {
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db}
}

func (r *appointmentRepository) FindAll(offset, limit int, search string, status string, startDate string, endDate string) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	query := r.db.Model(&models.Appointment{})

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Joins("LEFT JOIN patients ON patients.id = appointments.patient_id").
			Where("appointments.visit_number LIKE ? OR appointments.status LIKE ? OR appointments.complaint LIKE ? OR patients.name LIKE ?", searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if status != "" {
		query = query.Where("appointments.status = ?", status)
	}

	if startDate != "" {
		query = query.Where("appointments.appointment_date >= ?", startDate)
	}

	if endDate != "" {
		query = query.Where("appointments.appointment_date <= ?", endDate)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Patient").Preload("Patient.GenderData").Preload("Physiotherapist").Preload("ServiceMaster").Preload("TherapySession").Offset(offset).Limit(limit).Find(&appointments).Error
	return appointments, total, err
}

func (r *appointmentRepository) FindByID(id string) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Preload("Patient").Preload("Patient.GenderData").Preload("Physiotherapist").Preload("ServiceMaster").Preload("TherapySession").First(&appointment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) FindByPatientID(patientID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("patient_id = ?", patientID).Preload("Physiotherapist").Preload("ServiceMaster").Preload("TherapySession").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) FindByPhysiotherapistID(physioID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("physiotherapist_id = ?", physioID).Preload("Patient").Preload("Patient.GenderData").Preload("ServiceMaster").Preload("TherapySession").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *appointmentRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Appointment{}).Error
}

