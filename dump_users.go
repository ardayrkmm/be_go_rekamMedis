package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
)

func main() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "rekammedis-fisioterapi") // Wait, I don't know the project ID!
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	iter := client.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		fmt.Printf("Doc: %v\n", doc.Data())
	}
}
