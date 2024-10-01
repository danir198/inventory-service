// models/routes.go
package models

import (
	"net/http"

	middleware "github.com/danir198/inventory-service/auth"

	"github.com/danir198/inventory-service/handlers"

	"github.com/gorilla/mux"
)

type InventoryService interface {
	CheckAvailability(w http.ResponseWriter, r *http.Request)
	UpdateInventory(w http.ResponseWriter, r *http.Request)
	GetProduct(w http.ResponseWriter, r *http.Request)
	CreateProduct(w http.ResponseWriter, r *http.Request)
	DeleteProduct(w http.ResponseWriter, r *http.Request)
	ListProducts(w http.ResponseWriter, r *http.Request)
	SearchProducts(w http.ResponseWriter, r *http.Request)
}

type DbInventory struct {
	InventoryService InventoryService
}

func (s *DbInventory) InitializeRoutes() *mux.Router {
	router := mux.NewRouter()

	authHandler := &handlers.AuthHandler{}
	router.HandleFunc("/auth/token", authHandler.GenerateToken).Methods("POST")
	apiRouter := router.PathPrefix("/").Subrouter()
	apiRouter.Use(middleware.JWTAuth) // Apply the authentication middleware

	apiRouter.Handle("/products/{id}/availability", middleware.RequirePermission("check_availability", http.HandlerFunc(s.CheckAvailability))).Methods("GET")
	apiRouter.Handle("/products/{id}/inventory", middleware.RequirePermission("update_product", http.HandlerFunc(s.UpdateInventory))).Methods("PUT")
	apiRouter.Handle("/products/{id}", middleware.RequirePermission("view_product", http.HandlerFunc(s.GetProduct))).Methods("GET")
	apiRouter.Handle("/products", middleware.RequirePermission("add_product", http.HandlerFunc(s.CreateProduct))).Methods("POST")
	apiRouter.Handle("/products/{id}", middleware.RequirePermission("delete_product", http.HandlerFunc(s.DeleteProduct))).Methods("DELETE")
	apiRouter.Handle("/products", middleware.RequirePermission("view_all_products", http.HandlerFunc(s.ListProducts))).Methods("GET")
	apiRouter.Handle("/products/search", middleware.RequirePermission("search_products", http.HandlerFunc(s.SearchProducts))).Methods("GET")

	// apiRouter.HandleFunc("/products/{id}/availability", s.CheckAvailability).Methods("GET")
	// apiRouter.HandleFunc("/products/{id}/inventory", s.UpdateInventory).Methods("PUT")
	// apiRouter.HandleFunc("/products/{id}", s.GetProduct).Methods("GET")
	// apiRouter.HandleFunc("/products", s.CreateProduct).Methods("POST")
	// apiRouter.HandleFunc("/products/{id}", s.DeleteProduct).Methods("DELETE")
	// apiRouter.HandleFunc("/products", s.ListProducts).Methods("GET")

	// router.Use(middleware.BasicAuth) // Apply the authentication middleware
	// router.HandleFunc("/products/{id}/availability", s.CheckAvailability).Methods("GET")
	// router.HandleFunc("/products/{id}/inventory", s.UpdateInventory).Methods("PUT")
	// router.HandleFunc("/products/{id}", s.GetProduct).Methods("GET")
	// router.HandleFunc("/products", s.CreateProduct).Methods("POST")
	// router.HandleFunc("/products/{id}", s.DeleteProduct).Methods("DELETE")
	// router.HandleFunc("/products", s.ListProducts).Methods("GET")
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

func (s *DbInventory) SearchProducts(w http.ResponseWriter, r *http.Request) {
	s.InventoryService.SearchProducts(w, r)
}
