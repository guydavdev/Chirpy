package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fmt.Println("Starting server...")
	server := http.Server{Handler: mux, Addr: ":8080"}
	server.ListenAndServe()
}
