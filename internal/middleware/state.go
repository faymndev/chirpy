package middleware

import (
	"net/http"

	"github.com/faymndev/chirpy/internal/database"
)

type State struct {
	Db database.Queries
}

func (s *State) Middleware(next func(w http.ResponseWriter, r *http.Request, state *State)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next(w, r, s)
	})
}

/*
two potential usages:

1. middleware
mux.HandleFunc("/handle", state.Middleware(handle))
func handle(r, w, state)

2. attach to another struct (preferred?)
type Config struct {
	State *middleware.State
}

func (cfg *Config) handle() {
	we can access cfg.State inside the handler
}

mux.HandleFunc("/handle", handle)
*/
