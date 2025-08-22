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
	const filepathRoot = "."
	const port = "8080"
	apiConfig := apiConfig{}

	mux := http.NewServeMux()
	fileServerPrefix := "/app/"
	fileServerMux := http.StripPrefix(fileServerPrefix, http.FileServer(http.Dir(".")))
	mux.Handle(fileServerPrefix, apiConfig.middlewareMetricsInc(fileServerMux))

	mux.HandleFunc("GET /api/healthz", handleReadiness)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)

	server := http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving files from %s on port %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
