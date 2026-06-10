package routes

import (
	"encoding/json"
	"net/http"

	"github.com/faymndev/chirpy/internal/middleware"
)

func UseUsers(mux *http.ServeMux, state *middleware.State) {
	mux.Handle("POST /api/users", state.Middleware(handleCreateUser))
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		Email string `json:"email"`
	}

	// decode body
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	input := Input{}
	if err := decoder.Decode(&input); err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
	}

	user, err := s.Db.CreateUser(r.Context(), input.Email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
	}

	SendJSON(w, http.StatusCreated, user)
}
