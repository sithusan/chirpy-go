package main

import (
	"log"
	"net/http"
)

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
