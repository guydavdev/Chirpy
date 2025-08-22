package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/guydavdev/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const (
		filepathRoot     = "."
		port             = "8080"
		fileServerPrefix = "/app/"
	)

	godotenv.Load()
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	if err := dbConn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	defer dbConn.Close()

	apiConfig := apiConfig{db: database.New(dbConn)}

	mux := http.NewServeMux()
	fileServerMux := http.StripPrefix(fileServerPrefix, http.FileServer(http.Dir(".")))
	mux.Handle(fileServerPrefix, apiConfig.middlewareMetricsInc(fileServerMux))

	mux.HandleFunc("GET /api/healthz", handleReadiness)
	mux.HandleFunc("POST /api/validate_chirp", handlerChirpsValidate)
	mux.HandleFunc("POST /api/users", apiConfig.handlerUsers)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)

	server := http.Server{Handler: mux, Addr: ":" + port}

	log.Printf("Serving on port %s\n", port)
	log.Fatal(server.ListenAndServe())
}
