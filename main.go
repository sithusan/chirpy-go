package main

import (
	"log"
	"net/http"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	// static file server
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	mux.Handle("/healthz", http.HandlerFunc(healthzHandler))

	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
