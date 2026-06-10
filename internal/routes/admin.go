package routes

import (
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/faymndev/chirpy/internal/middleware"
)

type Metrics struct {
	State *middleware.State
	Hits  atomic.Int32
}

// apply metric routes
func UseAdmin(mux *http.ServeMux, state *middleware.State) *Metrics {
	metrics := &Metrics{State: state}
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
	if os.Getenv("PLATFORM") != "dev" {
		SendJSON(w, http.StatusForbidden, map[string]any{
			"error": "Cannot reset database outside of development",
		})
		return
	}

	cfg.Hits.Store(0)
	err := cfg.State.Db.Reset(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to reset database",
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: 0")
}
