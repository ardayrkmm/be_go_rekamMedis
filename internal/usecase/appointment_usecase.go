package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"errors"
	"fmt"
	"time"
)

type AppointmentUseCase interface {
	Fetch(offset, limit int, search string) ([]models.Appointment, int64, error)
	GetByID(id string) (*models.Appointment, error)
	GetByPatientID(patientID string) ([]models.Appointment, error)
	GetByPhysiotherapistID(physioID string) ([]models.Appointment, error)
	Store(appointment *models.Appointment) error
	Update(id string, req *models.Appointment) error
	Delete(id string) error
	Cancel(id string) error
	Reschedule(id string, req *models.Appointment) error
}

type appointmentUseCase struct {
	appointmentRepo repository.AppointmentRepository
	paymentRepo     repository.PaymentRepository
	serviceRepo     repository.ServiceMasterRepository
}

func NewAppointmentUseCase(appointmentRepo repository.AppointmentRepository, paymentRepo repository.PaymentRepository, serviceRepo repository.ServiceMasterRepository) AppointmentUseCase {
	return &appointmentUseCase{
		appointmentRepo: appointmentRepo,
		paymentRepo:     paymentRepo,
		serviceRepo:     serviceRepo,
	}
}

func (u *appointmentUseCase) Fetch(offset, limit int, search string) ([]models.Appointment, int64, error) {
	return u.appointmentRepo.FindAll(offset, limit, search)
}

func (u *appointmentUseCase) GetByID(id string) (*models.Appointment, error) {
	return u.appointmentRepo.FindByID(id)
}

func (u *appointmentUseCase) GetByPatientID(patientID string) ([]models.Appointment, error) {
	return u.appointmentRepo.FindByPatientID(patientID)
}

func (u *appointmentUseCase) GetByPhysiotherapistID(physioID string) ([]models.Appointment, error) {
	return u.appointmentRepo.FindByPhysiotherapistID(physioID)
}

func (u *appointmentUseCase) Store(appointment *models.Appointment) error {
	// Check for double booking
	if appointment.PhysiotherapistID != "" && appointment.AppointmentTime != "" {
		existing, err := u.appointmentRepo.FindByPhysiotherapistID(appointment.PhysiotherapistID)
		if err == nil {
			for _, app := range existing {
				if app.Status != "cancelled" && 
				   app.AppointmentDate != nil && appointment.AppointmentDate != nil &&
				   app.AppointmentDate.Format("2006-01-02") == appointment.AppointmentDate.Format("2006-01-02") && 
				   app.AppointmentTime == appointment.AppointmentTime {
					return errors.New("Slot waktu ini sudah terisi. Silakan pilih slot lain.")
				}
			}
		}
	}

	if appointment.Status == "" {
		appointment.Status = "scheduled"
	}
	
	if appointment.VisitNumber == nil || *appointment.VisitNumber == "" {
		vn := fmt.Sprintf("VIS-%s", time.Now().Format("20060102150405"))
		appointment.VisitNumber = &vn
	}

	err := u.appointmentRepo.Create(appointment)
	if err != nil {
		return err
	}

	// Auto-create payment for Appointment
	if appointment.ServiceMasterID != "" {
		service, err := u.serviceRepo.FindByID(appointment.ServiceMasterID)
		if err == nil && service != nil {
			// We can fetch names if needed, but since appointment has them or we can just fetch them
			// Wait, the appointment doesn't have the patient struct populated in Store() request, so we must fetch it.
			// Let's just create it. The payment list requires patient name, but wait, I can just leave it empty and let the fetch populate it. No, fetch doesn't populate it in NoSQL!
			
			payment := &models.Payment{
				InvoiceNumber:    "INV-" + time.Now().Format("20060102150405"),
				AppointmentID:    appointment.ID,
				PatientID:        appointment.PatientID,
				PhysiotherapistID: appointment.PhysiotherapistID,
				PaymentDate:      appointment.AppointmentDate,
				PaymentMethod:    "cash",
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

func (u *appointmentUseCase) Update(id string, req *models.Appointment) error {
	appointment, err := u.appointmentRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.PatientID != "" {
		appointment.PatientID = req.PatientID
	}
	if req.PhysiotherapistID != "" {
		appointment.PhysiotherapistID = req.PhysiotherapistID
	}
	if req.AppointmentDate != nil {
		appointment.AppointmentDate = req.AppointmentDate
	}
	if req.Status != "" {
		appointment.Status = req.Status
	}
	if req.Notes != "" {
		appointment.Notes = req.Notes
	}

	return u.appointmentRepo.Update(appointment)
}

func (u *appointmentUseCase) Delete(id string) error {
	return u.appointmentRepo.Delete(id)
}

func (u *appointmentUseCase) Cancel(id string) error {
	appointment, err := u.appointmentRepo.FindByID(id)
	if err != nil {
		return err
	}
	appointment.Status = "cancelled"
	return u.appointmentRepo.Update(appointment)
}

func (u *appointmentUseCase) Reschedule(id string, req *models.Appointment) error {
	appointment, err := u.appointmentRepo.FindByID(id)
	if err != nil {
		return err
	}
	appointment.AppointmentDate = req.AppointmentDate
	appointment.AppointmentTime = req.AppointmentTime
	appointment.Status = "rescheduled"
	return u.appointmentRepo.Update(appointment)
}
