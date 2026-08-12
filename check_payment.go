package main

import (
	"fmt"
	"log"
	"backend_go/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/rekam_medis?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var patientCount int64
	db.Model(&models.Patient{}).Where("id = ?", "4bc5bf18-e431-42f9-b8e0-54b1bb6278c8").Count(&patientCount)
	fmt.Printf("Patient exists in DB? %d\n", patientCount)
	
	var allPatients []models.Patient
	db.Find(&allPatients)
	fmt.Printf("Total patients in DB: %d\n", len(allPatients))
	for _, p := range allPatients {
		if p.ID == "4bc5bf18-e431-42f9-b8e0-54b1bb6278c8" {
			fmt.Printf("Found patient %s\n", p.Name)
		}
	}
}
