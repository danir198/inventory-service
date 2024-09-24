// models/handlers.go
package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InventoryServiceInterface interface {
	CheckAvailability(w http.ResponseWriter, r *http.Request)
	UpdateInventory(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
}

type InventoryService struct {
	// Define the fields and methods of the InventoryService type
	DB     *mongo.Database
	Logger *log.Logger
}

func (s *InventoryService) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Implement logic to check availability
	var result bson.M
	collection := s.DB.Collection("products")
	err := collection.FindOne(context.TODO(), bson.M{"product_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}

func (s *InventoryService) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Implement logic to update inventory

	var update bson.M
	err := json.NewDecoder(r.Body).Decode(&update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	collection := s.DB.Collection("products")
	filter := bson.M{"product_id": id}
	updateDoc := bson.M{"$set": update}
	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(context.TODO(), filter, updateDoc, opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *InventoryService) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Implement logic to get product
	var result bson.M
	collection := s.DB.Collection("products")
	err := collection.FindOne(context.TODO(), bson.M{"product_id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(result)
}
