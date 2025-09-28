package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func handlerChirpVerify(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type cleanedReturn struct {
		CleanedBody string `json:"cleaned_body"`
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
	} else {
		cleaned := cleanChirp(params.Body)
		respondJSON(w, http.StatusOK, cleanedReturn{
			CleanedBody: cleaned,
		})
	}
}

func respondJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}

func respondError(w http.ResponseWriter, errMsg string, code int) {
	type errorRetrun struct {
		Error string `json:"error"`
	}
	respBody := errorRetrun{
		Error: errMsg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
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
