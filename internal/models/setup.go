package models

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// Disable FK checks selama AutoMigrate untuk menghindari Error 3780
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS=1")

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

