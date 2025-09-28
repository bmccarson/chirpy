package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpVerify(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		invalidChirp(w, "Error decoding parameters", 500)
		return
	}

	if len(params.Body) > 140 {
		log.Print("Error: Body longer than 140 char")
		invalidChirp(w, "Body longer than 140 char", 400)
	} else {
		validChirp(w)
	}
}

func validChirp(w http.ResponseWriter) {
	type returnVals struct {
		Valid bool `json:"valid"`
	}
	respBody := returnVals{
		Valid: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}

func invalidChirp(w http.ResponseWriter, errMsg string, code int) {
	type returnVals struct {
		Error string `json:"error"`
	}
	respBody := returnVals{
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
