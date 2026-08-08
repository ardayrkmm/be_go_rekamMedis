package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"backend_go/internal/models"
	"backend_go/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	db := utils.ConnectDB()
	ctx := context.Background()
	now := time.Now()

	categories := []string{"Dewasa", "Pediatri", "Extra Treatment"}
	for _, cName := range categories {
		ref := db.Collection("service_categories").NewDoc()
		_, err := ref.Set(ctx, models.ServiceCategory{
			ID:        ref.ID,
			Name:      cName,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			log.Printf("Error creating category %s: %v\n", cName, err)
		}
	}

	services := []models.ServiceMaster{
		{Name: "1 Sesi Home Visit", Category: "Dewasa", Description: "Electrical Stimulation, Ultrasound, infrared, Terapi Latihan, Manual Terapi", BasePrice: 100000, Duration: 60},
		{Name: "Terapi Tumbuh Kembang", Category: "Pediatri", BasePrice: 75000, Duration: 60},
		{Name: "Neurosenso", Category: "Pediatri", BasePrice: 50000, Duration: 60},
		{Name: "Massage Bayi", Category: "Pediatri", BasePrice: 50000, Duration: 60},
		{Name: "Massage Anak", Category: "Pediatri", BasePrice: 75000, Duration: 60},
		{Name: "Terapi Dada", Category: "Pediatri", BasePrice: 50000, Duration: 60},
		{Name: "Nebulizer", Category: "Pediatri", BasePrice: 50000, Duration: 60},
		{Name: "Electrical Stimulation", Category: "Extra Treatment", BasePrice: 20000, Duration: 30},
		{Name: "Ultrasound", Category: "Extra Treatment", BasePrice: 20000, Duration: 30},
		{Name: "Infrared", Category: "Extra Treatment", BasePrice: 25000, Duration: 30},
		{Name: "Terapi Latihan", Category: "Extra Treatment", BasePrice: 25000, Duration: 30},
		{Name: "Manual Terapi", Category: "Extra Treatment", BasePrice: 25000, Duration: 30},
		{Name: "Kinesio Tapping", Category: "Extra Treatment", BasePrice: 20000, Duration: 30},
		{Name: "Masker Nebulizer", Category: "Extra Treatment", BasePrice: 20000, Duration: 30},
	}

	for i, svc := range services {
		svc.Code = fmt.Sprintf("SVC-%s-%d", now.Format("20060102150405"), i)
		svc.IsActive = true
		svc.CreatedAt = now
		svc.UpdatedAt = now
		
		ref := db.Collection("servicemasters").NewDoc()
		svc.ID = ref.ID
		_, err := ref.Set(ctx, svc)
		if err != nil {
			log.Printf("Error creating service %s: %v\n", svc.Name, err)
		}
	}

	fmt.Println("Seeding completed!")
}
