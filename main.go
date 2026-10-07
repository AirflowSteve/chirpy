package main

import (
	"net/http"
)

func handlerHealthz(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {
	handler := http.NewServeMux()
	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}
	handler.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	handler.HandleFunc("/healthz", handlerHealthz)
	server.ListenAndServe()
}
