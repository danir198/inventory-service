// models/routes.go
package models

import (
	"net/http"

	"inventory-service/handlers"

	"github.com/gorilla/mux"
)

type DbInventory struct {
	InventoryService handlers.InventoryService
}

func (s *DbInventory) InitializeRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}/availability", s.CheckAvailability).Methods("GET")
	router.HandleFunc("/products/{id}/inventory", s.UpdateInventory).Methods("PUT")
	router.HandleFunc("/products/{id}", s.GetProduct).Methods("GET")
	return router
}

func (s *DbInventory) CheckAvailability(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.CheckAvailability(w, r)
}

func (s *DbInventory) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.UpdateInventory(w, r)
}

func (s *DbInventory) GetProduct(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.GetProduct(w, r)
}
