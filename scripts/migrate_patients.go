package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
)

func main() {
	godotenv.Load(".env")
	
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore: %v\n", err)
	}
	defer client.Close()

	iter := client.Collection("patients").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		data := doc.Data()
		updated := false

		if catID, ok := data["patient_category_id"]; ok {
			switch v := catID.(type) {
			case int64:
				data["patient_category_id"] = fmt.Sprintf("%d", v)
				updated = true
			case float64:
				data["patient_category_id"] = fmt.Sprintf("%.0f", v)
				updated = true
			}
		}

		if genID, ok := data["gender_id"]; ok {
			switch v := genID.(type) {
			case int64:
				data["gender_id"] = fmt.Sprintf("%d", v)
				updated = true
			case float64:
				data["gender_id"] = fmt.Sprintf("%.0f", v)
				updated = true
			}
		}

		if updated {
			_, err := doc.Ref.Set(ctx, data)
			if err != nil {
				log.Printf("Failed to update doc %s: %v", doc.Ref.ID, err)
			} else {
				fmt.Printf("Updated doc %s\n", doc.Ref.ID)
			}
		}
	}
	fmt.Println("Migration complete.")
}
