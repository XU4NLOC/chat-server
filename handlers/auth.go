package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"chat-server/auth"
	"chat-server/config"
	"chat-server/models"

	"github.com/lib/pq"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "username and password required", http.StatusBadRequest)
			return
		}

		if err := models.CreateUser(req.Username, req.Password); err != nil {
			// Check specifically for unique constraint violation (Postgres error 23505)
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
				http.Error(w, "username already taken", http.StatusConflict)
				return
			}
			// For all other errors, log the real cause and return 500
			log.Printf("CreateUser error: %v", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Login(cfg config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req authRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByUsername(req.Username)
		if err != nil || !user.CheckPassword(req.Password) {
			// Always return the same error — don't reveal whether
			// the username or the password was wrong
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateToken(user.ID, user.Username, cfg.JWTSecret)
		if err != nil {
			http.Error(w, "could not generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
}

