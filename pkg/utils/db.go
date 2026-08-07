package utils

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

func ConnectDB() *firestore.Client {
	ctx := context.Background()

	// Use service account credentials if path is provided, otherwise default to application default credentials
	var app *firebase.App
	var err error

	credentialsFile := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsFile != "" {
		opt := option.WithCredentialsFile(credentialsFile)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Use Application Default Credentials
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}

	return client
}
