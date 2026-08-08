package main
import (
	"context"
	"fmt"
	"log"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"github.com/joho/godotenv"
	"os"
)
func main() {
	godotenv.Load()
	ctx := context.Background()
	opt := option.WithCredentialsFile("rekamGo.json")
	config := &firebase.Config{ProjectID: os.Getenv("FIREBASE_PROJECT_ID")}
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil { log.Fatal(err) }
	client, err := app.Firestore(ctx)
	if err != nil { log.Fatal(err) }
	defer client.Close()

	iter := client.Collection("medicalrecords").Limit(1).Documents(ctx)
	doc, err := iter.Next()
	if err == iterator.Done { fmt.Println("No documents"); return }
	if err != nil { log.Fatal(err) }

	fmt.Printf("Document data: %#v\n", doc.Data())
}
