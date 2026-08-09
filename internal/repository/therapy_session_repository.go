package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type TherapySessionRepository interface {
	FindAll(offset, limit int) ([]models.TherapySession, int64, error)
	FindByID(id string) (*models.TherapySession, error)
	FindByAppointmentID(appointmentID string) ([]models.TherapySession, error)
	Create(session *models.TherapySession) error
	Update(session *models.TherapySession) error
	Delete(id string) error
	GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error)
	HasTreatedPatient(physioID string, patientID string) (bool, error)
}

type therapySessionRepository struct {
	db *gorm.DB
}

func NewTherapySessionRepository(db *gorm.DB) TherapySessionRepository {
	return &therapySessionRepository{db}
}

func (r *therapySessionRepository) FindAll(offset, limit int) ([]models.TherapySession, int64, error) {
	var sessions []models.TherapySession
	var total int64

	err := r.db.Model(&models.TherapySession{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Appointment").Preload("Physiotherapist").Preload("Patient").Offset(offset).Limit(limit).Find(&sessions).Error
	return sessions, total, err
}

func (r *therapySessionRepository) FindByID(id string) (*models.TherapySession, error) {
	var session models.TherapySession
	err := r.db.Preload("Appointment").Preload("Physiotherapist").Preload("Patient").First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *therapySessionRepository) FindByAppointmentID(appointmentID string) ([]models.TherapySession, error) {
	var sessions []models.TherapySession
	err := r.db.Where("appointment_id = ?", appointmentID).Preload("Physiotherapist").Find(&sessions).Error
	return sessions, err
}

func (r *therapySessionRepository) Create(session *models.TherapySession) error {
	return r.db.Create(session).Error
}

func (r *therapySessionRepository) Update(session *models.TherapySession) error {
	return r.db.Save(session).Error
}

func (r *therapySessionRepository) Delete(id string) error {
	return r.db.Delete(&models.TherapySession{}, id).Error
}

func (r *therapySessionRepository) GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error) {
	var sessions []models.TherapySession
	err := r.db.Where("therapy_date >= ? AND therapy_date <= ?", startDate, endDate).Preload("Appointment").Preload("Patient").Preload("Physiotherapist").Find(&sessions).Error
	return sessions, err
}

func (r *therapySessionRepository) HasTreatedPatient(physioID string, patientID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.TherapySession{}).
		Where("physiotherapist_id = ? AND patient_id = ?", physioID, patientID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
