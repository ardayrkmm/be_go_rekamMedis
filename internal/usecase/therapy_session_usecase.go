package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"time"
)

type TherapySessionUseCase interface {
	Fetch(offset, limit int) ([]models.TherapySession, int64, error)
	GetByID(id string) (*models.TherapySession, error)
	GetByAppointmentID(appointmentID string) ([]models.TherapySession, error)
	Store(session *models.TherapySession) error
	Update(id string, session *models.TherapySession) error
	Delete(id string) error
	GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error)
}

type therapySessionUseCase struct {
	sessionRepo repository.TherapySessionRepository
	paymentRepo repository.PaymentRepository
	serviceRepo repository.ServiceMasterRepository
}

func NewTherapySessionUseCase(sessionRepo repository.TherapySessionRepository, paymentRepo repository.PaymentRepository, serviceRepo repository.ServiceMasterRepository) TherapySessionUseCase {
	return &therapySessionUseCase{
		sessionRepo: sessionRepo,
		paymentRepo: paymentRepo,
		serviceRepo: serviceRepo,
	}
}

func (u *therapySessionUseCase) Fetch(offset, limit int) ([]models.TherapySession, int64, error) {
	return u.sessionRepo.FindAll(offset, limit)
}

func (u *therapySessionUseCase) GetByID(id string) (*models.TherapySession, error) {
	return u.sessionRepo.FindByID(id)
}

func (u *therapySessionUseCase) GetByAppointmentID(appointmentID string) ([]models.TherapySession, error) {
	return u.sessionRepo.FindByAppointmentID(appointmentID)
}

func (u *therapySessionUseCase) Store(session *models.TherapySession) error {
	err := u.sessionRepo.Create(session)
	if err != nil {
		return err
	}

	// Auto-create payment
	if session.ServiceMasterID != "" {
		service, err := u.serviceRepo.FindByID(session.ServiceMasterID)
		if err == nil && service != nil {
			payment := &models.Payment{
				InvoiceNumber:    "INV-" + time.Now().Format("20060102150405"),
				TherapySessionID: session.ID,
				PatientID:        session.PatientID,
				PhysiotherapistID: session.PhysiotherapistID,
				PaymentDate:      session.TherapyDate,
				PaymentMethod:    "cash", // default
				Status:           "pending",
				Subtotal:         service.BasePrice,
				Discount:         0,
				Tax:              0,
				Total:            service.BasePrice,
				PaymentDetails: []models.PaymentDetail{
					{
						ServiceMasterID: service.ID,
						ItemName:        service.Name,
						Quantity:        1,
						Price:           service.BasePrice,
						Subtotal:        service.BasePrice,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
				},
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			}
			u.paymentRepo.Create(payment)
		}
	}
	return nil
}

func (u *therapySessionUseCase) Update(id string, req *models.TherapySession) error {
	session, err := u.sessionRepo.FindByID(id)
	if err != nil {
		return err
	}

	session.AppointmentID = req.AppointmentID
	session.PatientID = req.PatientID
	session.PhysiotherapistID = req.PhysiotherapistID
	session.ServiceMasterID = req.ServiceMasterID
	session.TherapyDate = req.TherapyDate
	session.Notes = req.Notes
	session.Status = req.Status
	session.Complaint = req.Complaint
	session.TreatmentGiven = req.TreatmentGiven

	return u.sessionRepo.Update(session)
}

func (u *therapySessionUseCase) Delete(id string) error {
	return u.sessionRepo.Delete(id)
}

func (u *therapySessionUseCase) GetWeeklySchedule(startDate, endDate string) ([]models.TherapySession, error) {
	return u.sessionRepo.GetWeeklySchedule(startDate, endDate)
}
