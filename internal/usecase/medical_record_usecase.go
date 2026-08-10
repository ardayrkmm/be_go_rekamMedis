package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"fmt"
	"time"
)

type MedicalRecordUseCase interface {
	Fetch(offset, limit int, search string) ([]models.MedicalRecord, int64, error)
	GetByID(id string) (*models.MedicalRecord, error)
	GetByPatientID(patientID string) ([]models.MedicalRecord, error)
	Store(record *models.MedicalRecord) error
	Update(id string, record *models.MedicalRecord) error
	Delete(id string) error
}

type medicalRecordUseCase struct {
	recordRepo  repository.MedicalRecordRepository
	patientRepo repository.PatientRepository
	physioRepo  repository.PhysiotherapistRepository
	serviceRepo repository.ServiceMasterRepository
	appointmentRepo repository.AppointmentRepository
}

func NewMedicalRecordUseCase(
	recordRepo repository.MedicalRecordRepository,
	patientRepo repository.PatientRepository,
	physioRepo repository.PhysiotherapistRepository,
	serviceRepo repository.ServiceMasterRepository,
	appointmentRepo repository.AppointmentRepository,
) MedicalRecordUseCase {
	return &medicalRecordUseCase{
		recordRepo:  recordRepo,
		patientRepo: patientRepo,
		physioRepo:  physioRepo,
		serviceRepo: serviceRepo,
		appointmentRepo: appointmentRepo,
	}
}

func (u *medicalRecordUseCase) Fetch(offset, limit int, search string) ([]models.MedicalRecord, int64, error) {
	records, total, err := u.recordRepo.FindAll(offset, limit, search)
	if err != nil {
		return nil, 0, err
	}
	
	for i := range records {
		u.populateRelations(&records[i])
	}
	return records, total, nil
}

func (u *medicalRecordUseCase) GetByID(id string) (*models.MedicalRecord, error) {
	record, err := u.recordRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	u.populateRelations(record)
	return record, nil
}

func (u *medicalRecordUseCase) GetByPatientID(patientID string) ([]models.MedicalRecord, error) {
	records, err := u.recordRepo.FindByPatientID(patientID)
	if err != nil {
		return nil, err
	}
	for i := range records {
		u.populateRelations(&records[i])
	}
	return records, nil
}

func (u *medicalRecordUseCase) populateRelations(record *models.MedicalRecord) {
	if record.VisitNumber == nil || *record.VisitNumber == "" {
		vn := fmt.Sprintf("VIS-%s", time.Now().Format("20060102150405"))
		record.VisitNumber = &vn
		// Save the generated visit number immediately to DB
		u.recordRepo.Update(record)
	}

	if record.PatientID != "" {
		if patient, err := u.patientRepo.FindByID(record.PatientID); err == nil {
			record.Patient = patient
		}
	}
	if record.PhysiotherapistID != "" {
		if physio, err := u.physioRepo.FindByID(record.PhysiotherapistID); err == nil {
			record.Physiotherapist = physio
		}
	}
	if record.ServiceID != nil && *record.ServiceID != "" {
		if svc, err := u.serviceRepo.FindByID(*record.ServiceID); err == nil {
			record.Service = svc
		}
	}
}

func (u *medicalRecordUseCase) Store(record *models.MedicalRecord) error {
	// If visit number is empty, try to get it from Appointment, or generate a new one
	if record.VisitNumber == nil || *record.VisitNumber == "" {
		if record.AppointmentID != nil && *record.AppointmentID != "" {
			if apt, err := u.appointmentRepo.FindByID(*record.AppointmentID); err == nil && apt.VisitNumber != nil {
				record.VisitNumber = apt.VisitNumber
			}
		}
		
		if record.VisitNumber == nil || *record.VisitNumber == "" {
			vn := fmt.Sprintf("VIS-%s", time.Now().Format("20060102150405"))
			record.VisitNumber = &vn
		}
	}
	return u.recordRepo.Create(record)
}

func (u *medicalRecordUseCase) Update(id string, req *models.MedicalRecord) error {
	record, err := u.recordRepo.FindByID(id)
	if err != nil {
		return err
	}

	if req.VisitNumber != nil && *req.VisitNumber != "" {
		record.VisitNumber = req.VisitNumber
	}
	record.PatientID = req.PatientID
	record.ServiceID = req.ServiceID
	record.PhysiotherapistID = req.PhysiotherapistID
	record.AppointmentID = req.AppointmentID
	record.ExaminationDate = req.ExaminationDate
	record.Anamnesis = req.Anamnesis
	record.Diagnosis = req.Diagnosis
	record.Therapy = req.Therapy
	record.Notes = req.Notes

	return u.recordRepo.Update(record)
}

func (u *medicalRecordUseCase) Delete(id string) error {
	return u.recordRepo.Delete(id)
}
