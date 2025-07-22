package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/clementine-tw/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	platform       string
	db             *database.Queries
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (a *apiConfig) handlerMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.Header().Add("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
	`, a.fileserverHits.Load())
	if err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func main() {

	// load environment variables
	godotenv.Load()

	// connect to database
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting database: %v", err)
	}
	dbQueries := database.New(db)

	// initialize config
	apiCfg := &apiConfig{
		platform: platform,
		db:       dbQueries,
	}

	const port = "8080"

	mux := http.NewServeMux()

	// healthcheck api
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// metrics api
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	// reset users
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetUsers)

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)

	// serve static files
	fileServer := apiCfg.middlewareMetricsInc(
		http.FileServer(http.Dir(".")),
	)
	mux.Handle("/app/", http.StripPrefix("/app/", fileServer))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port)
	server.ListenAndServe()
}
