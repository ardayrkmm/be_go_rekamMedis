package main

import (
	"context"
	"log"
	"os"

	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/firestore"
)

func main() {
	if err := godotenv.Load(`E:\Bisnis\RekamMedis\backend_go_firebase\.env`); err != nil {
		log.Println("Error loading .env file:", err)
	}
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", `E:\Bisnis\RekamMedis\backend_go_firebase\rekamGo.json`)

	db := utils.ConnectDB()
	ctx := context.Background()

	iter := db.Collection("users").Where("Role", "==", "admin").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		_, err = doc.Ref.Set(ctx, map[string]interface{}{
			"Name": "Pemilik",
		}, firestore.MergeAll)
		
		if err != nil {
			log.Printf("Failed to update doc %s: %v", doc.Ref.ID, err)
		} else {
			log.Printf("Updated admin %s name to Pemilik", doc.Ref.ID)
		}
	}
	
	log.Println("Done")
}
