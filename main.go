package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, _ *http.Request) {
	hits := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("Hits: %d", hits))
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, _ *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hits Reset")
}

func main() {
	const port = "8080"

	apiCfg := &apiConfig{}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serverMux.HandleFunc("/healthz", handlerReady)
	serverMux.HandleFunc("/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("/reset", apiCfg.handlerResetHits)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReady(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
