package routes

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	Hits atomic.Int32
}

// apply metric routes
func UseMetrics(mux *http.ServeMux) *Metrics {
	metrics := &Metrics{}
	mux.HandleFunc("GET /admin/metrics", metrics.handleMetrics)
	mux.HandleFunc("POST /admin/reset", metrics.handleReset)
	return metrics
}

func (cfg *Metrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.Hits.Add(1)
		next.ServeHTTP(w, r)
	})
}

const adminMetricsTemplate string = `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`

func (cfg *Metrics) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, adminMetricsTemplate, cfg.Hits.Load())
}

func (cfg *Metrics) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.Hits.Store(0)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: 0")
}
