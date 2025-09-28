package main

import (
	"fmt"
	"io"
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

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, hits))
}

func (cfg *apiConfig) handlerResetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Hits Reset")
}
