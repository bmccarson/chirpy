package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	apiCfg := &apiConfig{}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serverMux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets")))))
	serverMux.HandleFunc("GET /api/healthz", handlerReady)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetHits)
	serverMux.HandleFunc("POST /api/validate_chirp", handlerChirpVerify)

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReady(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
