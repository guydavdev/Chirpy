package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Println("Starting server...")
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}
