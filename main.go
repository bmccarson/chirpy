package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	serverMux := http.NewServeMux()
	serverMux.Handle("/", http.FileServer(http.Dir("")))

	server := &http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	fmt.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
