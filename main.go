package main

import (
	"fmt"
	"net/http"
)

type healthHandler struct{}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/healthz", healthHandler{})
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("public"))))
	server := &http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}
