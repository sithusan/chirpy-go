package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHit atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHit.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	res := fmt.Sprintf("Hits: %v", cfg.fileServerHit.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(res))
}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHit.Store(0)
	w.WriteHeader(http.StatusOK)
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	port := ":8080"
	mux := http.NewServeMux()

	// static file server
	apiCfg := &apiConfig{}
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	// health check
	mux.HandleFunc("/healthz", healthzHandler)

	// metrics
	mux.HandleFunc("/metrics", apiCfg.metricHandler)

	// reset
	mux.HandleFunc("/reset", apiCfg.resetHandler)

	// Serve the server
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
