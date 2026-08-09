package main

import (
	"encoding/json"
	"fmt"
	"backend_go/internal/models"
	"github.com/joho/godotenv"
	"os"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load()
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to DB:", err)
		return
	}

	var appointments []models.Appointment
	db.Preload("TherapySession").Find(&appointments)

	b, _ := json.MarshalIndent(appointments, "", "  ")
	fmt.Println(string(b))
}
