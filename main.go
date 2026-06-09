package main

import (
	"fmt"
	"net/http"

	"github.com/faymndev/chirpy/internal/middleware"
	"github.com/faymndev/chirpy/internal/routes"
)

func main() {
	mux := http.NewServeMux()

	metrics := routes.UseMetrics(mux)
	routes.UseChirp(mux)

	mux.Handle("GET /api/healthz", middleware.MiddlewareLog((healthHandler{})))
	mux.Handle("/app/", metrics.Middleware(http.StripPrefix("/app", http.FileServer(http.Dir("public")))))

	server := &http.Server{Addr: ":8080", Handler: mux}
	server.ListenAndServe()
}

type healthHandler struct{}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}
