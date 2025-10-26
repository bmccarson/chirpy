package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bmccarson/chirpy/internal/auth"
	"github.com/bmccarson/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondError(w, "Error decoding parameters", 500)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondError(w, "error hashing password", 500)
		return
	}

	usr, err := cfg.database.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondError(w, "Error creating user", 500)
		return
	}

	respondJSON(w, 201, User{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	})
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondError(w, "Forbidden", 403)
		return
	}
	err := cfg.database.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondError(w, "Error creating user", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "Users reset sucessful")
}

func (cfg *apiConfig) handlerLoginUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type user struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	req := request{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding request parameters: %s", err)
		respondError(w, "Error decoding request parameters", 500)
		return
	}

	usr, err := cfg.database.GetUserByEmail(context.Background(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, "user not found", 404)
			return
		}
	}

	match, err := auth.CheckPasswordHash(req.Password, usr.HashedPassword)
	if err != nil {
		log.Printf("error checking password hash: %v", err)
		return
	}

	if match {
		respondJSON(w, 200, user{
			ID:        usr.ID,
			CreatedAt: usr.CreatedAt,
			UpdatedAt: usr.UpdatedAt,
			Email:     usr.Email,
		})
	} else {
		respondJSON(w, 401, "")
	}
}
