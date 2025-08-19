package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServerPrefix := "/app/"
	fileServerMux := http.StripPrefix(fileServerPrefix, http.FileServer(http.Dir(".")))
	mux.Handle(fileServerPrefix, fileServerMux)
	mux.HandleFunc("GET /healthz", handleReadiness)

	fmt.Println("Starting server...")
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}

func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
