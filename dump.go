package main

import (
	"context"
	"fmt"

	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	db := utils.ConnectDB()
	ctx := context.Background()

	iter := db.Collection("users").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		fmt.Printf("Doc: %v\n", doc.Data())
	}
}
