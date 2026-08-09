package main

import (
	"context"
	"fmt"
	"log"

	"backend_go/pkg/utils"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	db := utils.ConnectDB()
	ctx := context.Background()

	iter := db.Collection("users").Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", doc.Data())
}
