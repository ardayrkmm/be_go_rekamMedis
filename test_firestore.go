package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("rekamGo.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	fmt.Println("--- APPOINTMENTS ---")
	iter1 := client.Collection("appointments").Limit(2).Documents(ctx)
	for {
		doc, err := iter1.Next()
		if err != nil {
			break
		}
		fmt.Printf("Doc ID: %s | Data: %v\n", doc.Ref.ID, doc.Data())
	}

	fmt.Println("--- THERAPY SESSIONS ---")
	iter2 := client.Collection("therapy_sessions").Limit(2).Documents(ctx)
	for {
		doc, err := iter2.Next()
		if err != nil {
			break
		}
		fmt.Printf("Doc ID: %s | Data: %v\n", doc.Ref.ID, doc.Data())
	}
}
