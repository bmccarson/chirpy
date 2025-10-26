package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bmccarson/chirpy/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"

	err := godotenv.Load()
	if err != nil {
		log.Fatal("could not load env")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Couldn't connect to database")
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		database: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}

	serverMux := http.NewServeMux()
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	serverMux.Handle("/app/assets/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/assets", http.FileServer(http.Dir("./assets")))))
	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLoginUser)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	serverMux.HandleFunc("GET /api/healthz", handlerReady)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)

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
	_, _ = io.WriteString(w, "OK")
}
