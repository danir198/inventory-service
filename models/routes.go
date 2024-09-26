// models/routes.go
package models

import (
	"net/http"

	"github.com/gorilla/mux"
)

type InventoryService interface {
	CheckAvailability(w http.ResponseWriter, r *http.Request)
	UpdateInventory(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
}

type DbInventory struct {
	InventoryService InventoryService
}

func (s *DbInventory) InitializeRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/products/{id}/availability", s.CheckAvailability).Methods("GET")
	router.HandleFunc("/products/{id}/inventory", s.UpdateInventory).Methods("PUT")
	router.HandleFunc("/products/{id}", s.GetProduct).Methods("GET")
	router.HandleFunc("/products", s.CreateProduct).Methods("POST")
	router.HandleFunc("/products/{id}", s.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/products", s.ListProducts).Methods("GET")
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
func (s *DbInventory) CreateProduct(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.CreateProduct(w, r)
}
func (s *DbInventory) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.DeleteProduct(w, r)
}
func (s *DbInventory) ListProducts(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.ListProducts(w, r)
}
