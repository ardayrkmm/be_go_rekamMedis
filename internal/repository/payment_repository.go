package repository

import (
	"backend_go/internal/models"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	FindAll(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error)
	FindByID(id string) (*models.Payment, error)
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
	FindByAppointmentID(appointmentID string) (*models.Payment, error)
	FindByTherapySessionID(therapySessionID string) (*models.Payment, error)
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) FindAll(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("payment_date >= ? AND payment_date <= ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("invoice_number LIKE ? OR patient_name LIKE ? OR physiotherapist_name LIKE ?", searchTerm, searchTerm, searchTerm)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Appointment").Preload("Patient").Preload("Patient.GenderData").Preload("PaymentDetails.ServiceMaster").Offset(offset).Limit(limit).Find(&payments).Error
	return payments, total, err
}

func (r *paymentRepository) FindByID(id string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Appointment").Preload("Patient").Preload("Patient.GenderData").Preload("PaymentDetails.ServiceMaster").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

func (r *paymentRepository) FindByAppointmentID(appointmentID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("appointment_id = ?", appointmentID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByTherapySessionID(therapySessionID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("therapy_session_id = ?", therapySessionID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}
