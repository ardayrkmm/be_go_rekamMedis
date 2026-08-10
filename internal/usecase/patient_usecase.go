package usecase

import (
	"backend_go/internal/models"
	"backend_go/internal/repository"
	"time"
)

type PatientUseCase interface {
	Fetch(offset, limit int, search string) ([]models.Patient, int64, error)
	GetByID(id string) (*models.Patient, error)
	Store(patient *models.Patient) error
	Update(id string, req *models.Patient) error
	Delete(id string) error
	Restore(id string) error
}

type patientUseCase struct {
	patientRepo repository.PatientRepository
}

func NewPatientUseCase(patientRepo repository.PatientRepository) PatientUseCase {
	return &patientUseCase{
		patientRepo: patientRepo,
	}
}

func (u *patientUseCase) Fetch(offset, limit int, search string) ([]models.Patient, int64, error) {
	return u.patientRepo.FindAll(offset, limit, search)
}

func (u *patientUseCase) GetByID(id string) (*models.Patient, error) {
	return u.patientRepo.FindByID(id)
}

func (u *patientUseCase) Store(patient *models.Patient) error {
	if patient.MedicalRecordNumber == "" {
		patient.MedicalRecordNumber = "RM-" + time.Now().Format("20060102150405")
	}
	return u.patientRepo.Create(patient)
}

func (u *patientUseCase) Update(id string, req *models.Patient) error {
	patient, err := u.patientRepo.FindByID(id)
	if err != nil {
		return err
	}

	// Update fields
	patient.Name = req.Name
	patient.Nik = req.Nik
	if req.MedicalRecordNumber != "" {
		patient.MedicalRecordNumber = req.MedicalRecordNumber
	} else if patient.MedicalRecordNumber == "" {
		patient.MedicalRecordNumber = "RM-" + time.Now().Format("20060102150405")
	}
	patient.BirthDate = req.BirthDate
	patient.PatientCategoryID = req.PatientCategoryID
	patient.GenderID = req.GenderID
	patient.BloodType = req.BloodType
	patient.Address = req.Address
	patient.Phone = req.Phone
	patient.Email = req.Email
	patient.Occupation = req.Occupation
	patient.MaritalStatus = req.MaritalStatus
	patient.EmergencyContactName = req.EmergencyContactName
	patient.EmergencyContactPhone = req.EmergencyContactPhone
	patient.MedicalHistory = req.MedicalHistory
	patient.Allergies = req.Allergies

	return u.patientRepo.Update(patient)
}

func (u *patientUseCase) Delete(id string) error {
	return u.patientRepo.Delete(id)
}

func (u *patientUseCase) Restore(id string) error {
	return u.patientRepo.Restore(id)
}
