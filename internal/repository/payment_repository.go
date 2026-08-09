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
		query = query.Where("payment_date >= ? AND payment_date <= ?", startDate, endDate)
	}

	// For search, we might need to join or just filter by invoice_number
	if search != "" {
		query = query.Where("invoice_number LIKE ?", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Appointment").Preload("Patient").Preload("PaymentDetails").Offset(offset).Limit(limit).Find(&payments).Error
	return payments, total, err
}

func (r *paymentRepository) FindByID(id string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Appointment").Preload("Patient").Preload("PaymentDetails").First(&payment, "id = ?", id).Error
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
