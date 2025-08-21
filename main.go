package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

func main() {
	apiConfig := apiConfig{}
	mux := http.NewServeMux()

	fileServerPrefix := "/app/"
	fileServerMux := http.StripPrefix(fileServerPrefix, http.FileServer(http.Dir(".")))
	mux.Handle(fileServerPrefix, apiConfig.middlewareMetricsInc(fileServerMux))
	mux.HandleFunc("/healthz", handleReadiness)
	mux.HandleFunc("/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("/reset", apiConfig.handleReset)

	fmt.Println("Starting server...")
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hits))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

type apiConfig struct {
	fileserverHits atomic.Int32
}
