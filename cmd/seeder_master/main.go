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
	if err := godotenv.Load("../../.env"); err != nil {
		godotenv.Load(".env")
	}

	db := utils.ConnectDB()
	ctx := context.Background()

	// 1. Patient Categories
	categories := []models.PatientCategory{
		{Name: "Umum", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "BPJS", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for i, cat := range categories {
		docID := []string{"1", "2"}[i]
		if _, err := db.Collection("patientcategories").Doc(docID).Set(ctx, cat); err != nil {
			log.Fatalf("Failed to create category: %v", err)
		}
	}

	// 2. Genders
	genders := []models.Gender{
		{Name: "Laki-laki", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{Name: "Perempuan", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for i, g := range genders {
		docID := []string{"1", "2"}[i]
		if _, err := db.Collection("genders").Doc(docID).Set(ctx, g); err != nil {
			log.Fatalf("Failed to create gender: %v", err)
		}
	}

	log.Println("Master seeder completed successfully!")
}
