package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"errors"
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
	HasTreatedPatient(physioID, patientID string) (bool, error)
}

type therapySessionUseCase struct {
	sessionRepo repository.TherapySessionRepository
	paymentRepo repository.PaymentRepository
	serviceRepo repository.ServiceMasterRepository
	recordRepo  repository.MedicalRecordRepository
}

func NewTherapySessionUseCase(sessionRepo repository.TherapySessionRepository, paymentRepo repository.PaymentRepository, serviceRepo repository.ServiceMasterRepository, recordRepo repository.MedicalRecordRepository) TherapySessionUseCase {
	return &therapySessionUseCase{
		sessionRepo: sessionRepo,
		paymentRepo: paymentRepo,
		serviceRepo: serviceRepo,
		recordRepo:  recordRepo,
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
	var serviceIDs []string
	if len(session.ServiceMasterIDs) > 0 {
		serviceIDs = session.ServiceMasterIDs
	} else if session.ServiceMasterID != "" {
		serviceIDs = []string{session.ServiceMasterID}
	}

	if len(serviceIDs) > 0 {
		var details []models.PaymentDetail
		var subtotal float64
		for _, id := range serviceIDs {
			service, err := u.serviceRepo.FindByID(id)
			if err == nil && service != nil {
				details = append(details, models.PaymentDetail{
					ServiceMasterID: service.ID,
					ItemName:        service.Name,
					Quantity:        1,
					Price:           service.BasePrice,
					Subtotal:        service.BasePrice,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				})
				subtotal += service.BasePrice
			}
		}

		if len(details) > 0 {
			payment := &models.Payment{
				InvoiceNumber:     "INV-" + time.Now().Format("20060102150405"),
				TherapySessionID:  session.ID,
				PatientID:         session.PatientID,
				PhysiotherapistID: session.PhysiotherapistID,
				PaymentDate:       session.TherapyDate,
				PaymentMethod:     "cash", // default
				Status:            "pending",
				Subtotal:          subtotal,
				Discount:          0,
				Tax:               0,
				Total:             subtotal,
				PaymentDetails:    details,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
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

	// Validasi: Sesi yang sudah selesai tidak boleh kembali ke scheduled
	if session.Status == "completed" && req.Status == "scheduled" {
		return errors.New("Sesi yang sudah diselesaikan tidak dapat diubah kembali menjadi Scheduled")
	}

	// Validasi: Untuk menyelesaikan sesi (status = completed), harus ada Data Klinis (Medical Record)
	if req.Status == "completed" && session.Status != "completed" {
		records, err := u.recordRepo.FindByPatientID(session.PatientID)
		hasRecord := false
		if err == nil {
			for _, rec := range records {
				if rec.AppointmentID != nil && *rec.AppointmentID == session.AppointmentID {
					hasRecord = true
					break
				}
			}
		}
		if !hasRecord {
			return errors.New("Tidak dapat menyelesaikan sesi: Data klinis (Rekam Medis) wajib diisi terlebih dahulu")
		}
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

func (u *therapySessionUseCase) HasTreatedPatient(physioID, patientID string) (bool, error) {
	return u.sessionRepo.HasTreatedPatient(physioID, patientID)
}
