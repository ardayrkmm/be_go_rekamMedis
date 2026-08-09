package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
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
	db *gorm.DB
}

func NewAppointmentRepository(db *gorm.DB) AppointmentRepository {
	return &appointmentRepository{db}
}

func (r *appointmentRepository) FindAll(offset, limit int) ([]models.Appointment, int64, error) {
	var appointments []models.Appointment
	var total int64

	err := r.db.Model(&models.Appointment{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Patient").Preload("Physiotherapist").Preload("ServiceMaster").Offset(offset).Limit(limit).Find(&appointments).Error
	return appointments, total, err
}

func (r *appointmentRepository) FindByID(id string) (*models.Appointment, error) {
	var appointment models.Appointment
	err := r.db.Preload("Patient").Preload("Physiotherapist").Preload("ServiceMaster").First(&appointment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &appointment, nil
}

func (r *appointmentRepository) FindByPatientID(patientID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("patient_id = ?", patientID).Preload("Physiotherapist").Preload("ServiceMaster").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) FindByPhysiotherapistID(physioID string) ([]models.Appointment, error) {
	var appointments []models.Appointment
	err := r.db.Where("physiotherapist_id = ?", physioID).Preload("Patient").Preload("ServiceMaster").Find(&appointments).Error
	return appointments, err
}

func (r *appointmentRepository) Create(appointment *models.Appointment) error {
	return r.db.Create(appointment).Error
}

func (r *appointmentRepository) Update(appointment *models.Appointment) error {
	return r.db.Save(appointment).Error
}

func (r *appointmentRepository) Delete(id string) error {
	return r.db.Delete(&models.Appointment{}, id).Error
}
