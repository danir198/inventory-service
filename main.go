// Step 7: Set up the main entry point
// main.go
package main

import (
	"context"
	"inventory-service/handlers"
	"inventory-service/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Check if MONGO_URI is set
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("required key MONGO_URI missing value")
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: %v", err)
	}

	db, err := models.InitializeDatabase(client, os.Getenv("DATABASE_NAME"))
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize logger
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	service := &handlers.InventoryService{
		DB:     db,
		Logger: logger,
	}

	dbInventory := &models.DbInventory{
		InventoryService: *service,
	}

	log.Printf("Server started")

	router := dbInventory.InitializeRoutes()
	log.Fatal(http.ListenAndServe(os.Getenv("SERVER_ADDRESS"), router))
}
