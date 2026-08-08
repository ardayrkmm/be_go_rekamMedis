package main

import (
	"context"
	"log"
	"os"

	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func main() {
	if err := godotenv.Load(`E:\Bisnis\RekamMedis\backend_go_firebase\.env`); err != nil {
		log.Println("Error loading .env file:", err)
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", `E:\Bisnis\RekamMedis\backend_go_firebase\rekamGo.json`)

	db := utils.ConnectDB()
	ctx := context.Background()

	iter := db.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		var data map[string]interface{}
		doc.DataTo(&data)
		log.Printf("User: %v", data)
	}
	
	log.Println("Done")
}
