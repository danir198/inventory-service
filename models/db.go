// models/db.go
package models

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeDatabase(client *mongo.Client, dbName string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if the database exists
	db := client.Database(dbName)
	collections, err := db.ListCollectionNames(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Create collections if they do not exist
	requiredCollections := []string{"products"}
	for _, collection := range requiredCollections {
		if !contains(collections, collection) {
			err := db.CreateCollection(ctx, collection)
			if err != nil {
				return nil, err
			}
			log.Printf("Created collection: %s", collection)
		}
	}

	// Create index on product_id field in the products collection
	collection := db.Collection("products")
	_, err = collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.M{"product_id": 1},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		// Check if the error is because the index already exists
		if !strings.Contains(err.Error(), "already exists") {
			return nil, fmt.Errorf("error creating index: %w", err)
		}
		log.Println("Index on product_id already exists")
	} else {
		log.Println("Index on product_id created successfully")
	}

	return db, nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
