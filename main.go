// Step 7: Set up the main entry point
// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/danir198/inventory-service/db"

	"github.com/danir198/inventory-service/handlers"
	"github.com/danir198/inventory-service/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Load environment variables from .env file
	if err := loadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	mongoURI, databaseName, serverAddress := getEnvVars()
	client, err := createMongoClient(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	db := db.NewDatabase(client, databaseName)

	if err := db.InitializeDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// db, err := models.InitializeDatabase(client, databaseName)
	// if err != nil {
	// 	log.Fatalf("Failed to initialize database: %v", err)
	// }

	// Initialize logger
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	service := &handlers.InventoryService{
		DB:     db,
		Logger: logger,
	}

	dbInventory := &models.DbInventory{
		InventoryService: service,
	}

	log.Printf("Server started")

	router := dbInventory.InitializeRoutes()
	log.Fatal(http.ListenAndServe(serverAddress, router))
}

func loadEnv() error {
	return godotenv.Load()
}

func getEnvVars() (string, string, string) {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatalf("required key MONGO_URI missing value")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	if databaseName == "" {
		log.Fatalf("required key DATABASE_NAME missing value")
	}

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		log.Fatalf("required key SERVER_ADDRESS missing value")
	}

	return mongoURI, databaseName, serverAddress
}

func createMongoClient(mongoURI string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	return mongo.NewClient(clientOptions)
}
