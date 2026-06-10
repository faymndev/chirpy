package routes

import (
	"encoding/json"
	"net/http"

	"github.com/faymndev/chirpy/internal/database"
)

func UseUsers(mux *http.ServeMux, db *database.Queries) {
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
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

		user, err := db.CreateUser(r.Context(), input.Email)
		if err != nil {
			SendJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "Something went wrong",
			})
		}

		SendJSON(w, http.StatusCreated, user)
	})
}
