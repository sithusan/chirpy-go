package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sithusan/chirpy-go/internal/database"
)

type apiConfig struct {
	db            *database.Queries
	fileServerHit atomic.Int32
}

func initiateDB() *database.Queries {

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Connection Problem: %v", err)
	}

	return database.New(dbConn)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Closure, needs to run in every request that comes to file server.
		cfg.fileServerHit.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	res := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileServerHit.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
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

	apiCfg := &apiConfig{
		fileServerHit: atomic.Int32{},
		db:            initiateDB(),
	}

	// static file server
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	// health check
	mux.HandleFunc("GET /api/healthz", healthzHandler)

	// metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)

	// reset
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

	// Serve the server
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
