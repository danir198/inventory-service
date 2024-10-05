package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	middleware "github.com/danir198/inventory-service/auth"
)

type AuthHandler struct{}
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (a *AuthHandler) GenerateToken(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Split the API_USERNAME environment variable into a slice
	validUsernames := strings.Split(os.Getenv("API_USERNAME"), ",")

	isValidUsername := false

	for _, validUsername := range validUsernames {
		if creds.Username == validUsername {
			isValidUsername = true
			break
		}
	}

	if !isValidUsername || creds.Password != os.Getenv("API_PASSWORD") {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	role := middleware.User
	if creds.Username == "admin" {
		role = middleware.Admin
	}

	token, err := middleware.GenerateToken(creds.Username, string(role))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
