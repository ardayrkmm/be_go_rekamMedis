package main

import (
	"log"
	"time"

	"backend_go/internal/models"
	"backend_go/pkg/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}

	db := utils.ConnectDB()

	// Reset semua tabel
	log.Println("Menghapus semua data...")
	tables := []string{
		"activity_logs", "notifications", "payment_details", "payments",
		"exercise_programs", "pain_assessments", "therapy_sessions",
		"medical_records", "appointments", "patients", "physiotherapists",
		"service_masters", "users", "genders", "patient_categories",
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=0")
	for _, t := range tables {
		db.Exec("TRUNCATE TABLE " + t)
	}
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
	log.Println("Semua tabel berhasil direset.")

	// Buat akun Admin (Pemilik)
	hashedPwd, err := utils.HashPassword("admin1234")
	if err != nil {
		log.Fatalf("Gagal hash password: %v", err)
	}

	admin := models.User{
		ID:        uuid.New().String(),
		Name:      "Pemilik",
		Email:     "arummyfisioterapi@pemilik.com",
		Password:  hashedPwd,
		Role:      string(models.RoleAdmin),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&admin).Error; err != nil {
		log.Fatalf("Gagal membuat akun admin: %v", err)
	}

	// Insert Genders
	genders := []models.Gender{
		{ID: uuid.New().String(), Name: "Laki-laki"},
		{ID: uuid.New().String(), Name: "Perempuan"},
	}
	for _, g := range genders {
		db.Create(&g)
	}

	// Insert Patient Categories
	categories := []models.PatientCategory{
		{ID: uuid.New().String(), Name: "Umum"},
		{ID: uuid.New().String(), Name: "Member"},
	}
	for _, c := range categories {
		db.Create(&c)
	}

	log.Println("=================================")
	log.Println("Seeder selesai!")
	log.Println("--- AKUN ADMIN ---")
	log.Printf("Email   : %s\n", admin.Email)
	log.Printf("Password: admin1234\n")
	log.Printf("Role    : %s\n", admin.Role)
	log.Println("=================================")
}
