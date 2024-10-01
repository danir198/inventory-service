// models/handlers.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/danir198/inventory-service/db"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryService struct {
	// Define the fields and methods of the InventoryService type
	DB     *db.Database
	Logger *log.Logger
}

func (s *InventoryService) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// Implement logic to check availability
	result, err := s.DB.FindProductByID(id)

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

	result, err := s.DB.UpdateProductByID(id, update)
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
	result, err := s.DB.FindProductByID(id)
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

func (s *InventoryService) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Add logic to insert the product into the database
	productMap, err := bson.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var productBson bson.M
	err = bson.Unmarshal(productMap, &productBson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	results, err := s.DB.CreateProduct(productBson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(results)
}

func (s *InventoryService) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	result, err := s.DB.DeleteProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *InventoryService) ListProducts(w http.ResponseWriter, r *http.Request) {
	results, err := s.DB.ListProducts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(results)
}

func (s *InventoryService) SearchProducts(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	results, err := s.DB.SearchProducts(filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(results)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
