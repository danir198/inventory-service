package middleware

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func BasicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if os.Getenv("ENABLE_AUTH") != "true" {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(auth, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(authParts[1])
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 || !validateCredentials(pair[0], pair[1]) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateCredentials(username, password string) bool {
	return username == os.Getenv("API_USERNAME") && password == os.Getenv("API_PASSWORD")
}

type contextKey string

const (
	//userContextKey contextKey = "user"
	ContextKeyUserRole contextKey = "userRole"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			http.Error(w, "Invalid role", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserRole, Role(role))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateToken(username,role string) (string, error) {
	expirationSeconds, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION"))
	if err != nil {
		return "", err
	}
	expirationTime := time.Now().Add(time.Duration(expirationSeconds) * time.Second)
	claims := &jwt.MapClaims{
		"user": username,
		"role": role,
		"exp":  expirationTime.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func RequirePermission(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextKeyUserRole).(Role)
		if !ok||!HasPermission(role, permission) {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
