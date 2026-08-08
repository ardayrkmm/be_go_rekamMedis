package main

import (
	"context"
	"fmt"
	"log"

	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
	"google.golang.org/api/iterator"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found")
	}

	db := utils.ConnectDB()
	ctx := context.Background()

	// Get all documents in the payments collection
	iter := db.Collection("payments").Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate over payments: %v", err)
		}

		// Delete the document
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("Failed to delete document %s: %v", doc.Ref.ID, err)
		} else {
			count++
		}
	}

	fmt.Printf("Successfully deleted %d payment records.\n", count)
}
