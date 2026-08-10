package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"errors"
)

type PaymentUseCase interface {
	Fetch(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error)
	GetByID(id string) (*models.Payment, error)
	Store(payment *models.Payment) error
	Update(id string, payment *models.Payment) error
	UpdateStatus(id string, status string) error
}

type paymentUseCase struct {
	paymentRepo repository.PaymentRepository
	patientRepo repository.PatientRepository
	physioRepo  repository.PhysiotherapistRepository
}

func NewPaymentUseCase(paymentRepo repository.PaymentRepository, patientRepo repository.PatientRepository, physioRepo repository.PhysiotherapistRepository) PaymentUseCase {
	return &paymentUseCase{
		paymentRepo: paymentRepo,
		patientRepo: patientRepo,
		physioRepo:  physioRepo,
	}
}

func (u *paymentUseCase) populateNames(payment *models.Payment) {
	if payment.PatientName == "" && payment.PatientID != "" {
		patient, err := u.patientRepo.FindByID(payment.PatientID)
		if err == nil && patient != nil {
			payment.PatientName = patient.Name
		}
	}
	if payment.PhysiotherapistName == "" && payment.PhysiotherapistID != "" {
		physio, err := u.physioRepo.FindByID(payment.PhysiotherapistID)
		if err == nil && physio != nil {
			payment.PhysiotherapistName = physio.Name
		}
	}
}

func (u *paymentUseCase) Fetch(offset, limit int, search, status, startDate, endDate string) ([]models.Payment, int64, error) {
	payments, total, err := u.paymentRepo.FindAll(offset, limit, search, status, startDate, endDate)
	if err == nil {
		for i := range payments {
			u.populateNames(&payments[i])
		}
	}
	return payments, total, err
}

func (u *paymentUseCase) GetByID(id string) (*models.Payment, error) {
	payment, err := u.paymentRepo.FindByID(id)
	if err == nil && payment != nil {
		u.populateNames(payment)
	}
	return payment, err
}

func (u *paymentUseCase) Store(payment *models.Payment) error {
	// Prevent duplicate payments for the same Therapy Session or Appointment
	if payment.TherapySessionID != "" {
		existing, err := u.paymentRepo.FindByTherapySessionID(payment.TherapySessionID)
		if err == nil && existing != nil {
			return errors.New("Pembayaran untuk sesi terapi ini sudah ada")
		}
	} else if payment.AppointmentID != "" {
		existing, err := u.paymentRepo.FindByAppointmentID(payment.AppointmentID)
		if err == nil && existing != nil {
			return errors.New("Pembayaran untuk janji ini sudah ada")
		}
	}

	u.populateNames(payment)
	return u.paymentRepo.Create(payment)
}

func (u *paymentUseCase) UpdateStatus(id string, status string) error {
	payment, err := u.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}

	payment.Status = status
	return u.paymentRepo.Update(payment)
}

func (u *paymentUseCase) Update(id string, req *models.Payment) error {
	existing, err := u.paymentRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Update only the fields that can be edited from form
	existing.PatientID = req.PatientID
	existing.PhysiotherapistID = req.PhysiotherapistID
	existing.TherapySessionID = req.TherapySessionID
	existing.AppointmentID = req.AppointmentID
	existing.Status = req.Status
	existing.PaymentMethod = req.PaymentMethod
	existing.Notes = req.Notes
	existing.Subtotal = req.Subtotal
	existing.Discount = req.Discount
	existing.Tax = req.Tax
	existing.Total = req.Total
	existing.PaymentDetails = req.PaymentDetails

	// Fetch names for updated patient and physio
	if existing.PatientID != "" {
		if p, err := u.patientRepo.FindByID(existing.PatientID); err == nil {
			existing.PatientName = p.Name
		}
	}
	if existing.PhysiotherapistID != "" {
		if p, err := u.physioRepo.FindByID(existing.PhysiotherapistID); err == nil {
			existing.PhysiotherapistName = p.Name
		}
	}

	return u.paymentRepo.Update(existing)
}
