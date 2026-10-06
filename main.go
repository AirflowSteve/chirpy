package main

import (
	"net/http"
)

func main() {
	handler := http.NewServeMux()
	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	handler.Handle("/", http.FileServer(http.Dir(".")))
	server.ListenAndServe()
}
