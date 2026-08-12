package models

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// Disable FK checks selama AutoMigrate untuk menghindari Error 3780
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS=1")

	// Drop old unique indexes to allow duplicate 0 values
	if db.Migrator().HasIndex(&Patient{}, "idx_patients_medical_record_number") {
		db.Migrator().DropIndex(&Patient{}, "idx_patients_medical_record_number")
	}
	if db.Migrator().HasIndex(&Patient{}, "idx_patients_nik") {
		db.Migrator().DropIndex(&Patient{}, "idx_patients_nik")
	}
	if db.Migrator().HasIndex(&Patient{}, "idx_patients_email") {
		db.Migrator().DropIndex(&Patient{}, "idx_patients_email")
	}
	if db.Migrator().HasIndex(&Physiotherapist{}, "idx_physiotherapists_sip") {
		db.Migrator().DropIndex(&Physiotherapist{}, "idx_physiotherapists_sip")
	}

	return db.AutoMigrate(
		&User{},
		&JwtBlocklist{},
		&PasswordResetToken{},
		&Gender{},
		&ServiceCategory{},
		&ServiceMaster{},
		&Physiotherapist{},
		&PatientCategory{},
		&Patient{},
		&MedicalRecord{},
		&Appointment{},
		&TherapySession{},
		&Payment{},
		&PaymentDetail{},
		&Notification{},
		&ActivityLog{},
		&PainAssessment{},
		&ExerciseProgram{},
	)
}

