package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
	"github.com/faymndev/chirpy/internal/routes"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	dbQueries := database.New(db)
	state := &middleware.State{Db: *dbQueries}
	mux := http.NewServeMux()

	metrics := routes.UseAdmin(mux, state)
	routes.UseChirp(mux, state)
	routes.UseUsers(mux, state)

	mux.Handle("GET /api/healthz", middleware.Log((healthHandler{})))
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
