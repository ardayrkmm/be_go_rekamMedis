package main

import (
	"context"
	"log"
	"time"

	"backend_go/internal/models"
	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}

	db := utils.ConnectDB()
	ctx := context.Background()

	// 1. Create Admin
	adminPassword, _ := utils.HashPassword("admin123")
	admin := models.User{
		Name:      "Pemilik",
		Email:     "admin@klinik.com",
		Password:  adminPassword,
		Role:      string(models.RoleAdmin),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Create Fisioterapis
	fisioPassword, _ := utils.HashPassword("fisio123")
	fisio := models.User{
		Name:      "Budi Fisioterapis",
		Email:     "budi@klinik.com",
		Password:  fisioPassword,
		Role:      string(models.RoleFisioterapis),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Write to Firestore
	adminRef := db.Collection("users").NewDoc()
	admin.ID = adminRef.ID
	if _, err := adminRef.Set(ctx, admin); err != nil {
		log.Fatalf("Failed to create admin: %v", err)
	}

	fisioRef := db.Collection("users").NewDoc()
	fisio.ID = fisioRef.ID
	if _, err := fisioRef.Set(ctx, fisio); err != nil {
		log.Fatalf("Failed to create fisioterapis: %v", err)
	}

	log.Println("Seeder completed successfully!")
	log.Println("--- AKUN LOGIN ---")
	log.Println("Admin: admin@klinik.com | Pass: admin123")
	log.Println("Fisio: budi@klinik.com  | Pass: fisio123")
}
