// models/models_test.go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCheckAvailability(t *testing.T) {
	req, err := http.NewRequest("GET", "/products/123/availability", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc((&InventoryService{}).CheckAvailability)

	router := mux.NewRouter()
	router.HandleFunc("/products/{id}/availability", handler)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
