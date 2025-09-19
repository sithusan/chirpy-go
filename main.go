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
	platform      string
}

func initiateDB() *database.Queries {
	dbURL := getEnvOrFail("DB_URL")
	dbConn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("Connection Problem: %v", err)
	}

	return database.New(dbConn)
}

func getEnvOrFail(key string) string {
	envVal := os.Getenv(key)

	if envVal == "" {
		log.Fatalf("%s must be set", key)
	}

	return envVal
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

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	godotenv.Load()

	apiCfg := &apiConfig{
		fileServerHit: atomic.Int32{},
		db:            initiateDB(),
		platform:      getEnvOrFail("PLATFORM"),
	}

	port := ":8080"
	mux := http.NewServeMux()

	// static file server
	fileServerHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileServerHandler))

	// health check
	mux.HandleFunc("GET /api/healthz", healthzHandler)

	// metrics
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricHandler)

	// reset
	mux.HandleFunc("POST /admin/reset", apiCfg.resetUserHandler)

	// chirps
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)

	// users
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	// Serve the server
	server := &http.Server{
		Addr:    port,
		Handler: mux,
	}

	log.Printf("Serving on port %v\n", port)
	log.Fatal(server.ListenAndServe())
}
