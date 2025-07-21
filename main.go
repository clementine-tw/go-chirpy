package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
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

func (a *apiConfig) handlerResetMetrics(w http.ResponseWriter, _ *http.Request) {
	a.fileserverHits.Store(0)
	w.Header().Add("Cache-Control", "no-cache")
	w.WriteHeader(http.StatusOK)
}

func main() {

	apiCfg := &apiConfig{}

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
	// reset metrics api
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerResetMetrics)

	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

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
