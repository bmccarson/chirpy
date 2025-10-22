package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/bmccarson/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string    `json:"body"`
		User uuid.UUID `json:"user_id"`
	}
	type newChirp struct {
		ChirpID   uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		User      uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondError(w, "Error decoding parameters", 500)
		return
	}

	if len(params.Body) > 140 {
		log.Print("Error: Body longer than 140 char")
		respondError(w, "Body longer than 140 char", 400)
		return
	}
	cleaned := cleanChirp(params.Body)

	chirp, err := cfg.database.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.User,
	})
	if err != nil {
		log.Printf("error creating chirp: %v", err)
		respondError(w, "database failed to create chirp", 500)
	}

	respondJSON(w, 201, newChirp{
		ChirpID:   chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User:      chirp.UserID,
	})
}

func cleanChirp(body string) string {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")

	return cleaned
}
