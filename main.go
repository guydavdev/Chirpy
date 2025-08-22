package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const (
		filepathRoot     = "."
		port             = "8080"
		fileServerPrefix = "/app/"
	)

	apiConfig := apiConfig{}

	mux := http.NewServeMux()
	fileServerMux := http.StripPrefix(fileServerPrefix, http.FileServer(http.Dir(".")))
	mux.Handle(fileServerPrefix, apiConfig.middlewareMetricsInc(fileServerMux))

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)

	server := http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
