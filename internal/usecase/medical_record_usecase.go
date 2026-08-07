package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
)

type MedicalRecordUseCase interface {
	Fetch(offset, limit int) ([]models.MedicalRecord, int64, error)
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
}

func NewMedicalRecordUseCase(
	recordRepo repository.MedicalRecordRepository,
	patientRepo repository.PatientRepository,
	physioRepo repository.PhysiotherapistRepository,
	serviceRepo repository.ServiceMasterRepository,
) MedicalRecordUseCase {
	return &medicalRecordUseCase{
		recordRepo:  recordRepo,
		patientRepo: patientRepo,
		physioRepo:  physioRepo,
		serviceRepo: serviceRepo,
	}
}

func (u *medicalRecordUseCase) Fetch(offset, limit int) ([]models.MedicalRecord, int64, error) {
	records, total, err := u.recordRepo.FindAll(offset, limit)
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
	return u.recordRepo.Create(record)
}

func (u *medicalRecordUseCase) Update(id string, req *models.MedicalRecord) error {
	record, err := u.recordRepo.FindByID(id)
	if err != nil {
		return err
	}

	record.VisitNumber = req.VisitNumber
	record.PatientID = req.PatientID
	record.ServiceID = req.ServiceID
	record.PhysiotherapistID = req.PhysiotherapistID
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
