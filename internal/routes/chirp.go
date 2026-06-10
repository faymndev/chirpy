package routes

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
	"github.com/google/uuid"
)

func UseChirp(mux *http.ServeMux, state *middleware.State) {
	mux.Handle("POST /api/chirps", state.Middleware(handleChirp))
}

func handleChirp(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		UserID uuid.UUID `json:"user_id"`
		Body   string    `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	input := Input{}
	err := decoder.Decode(&input)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	} else if utf8.RuneCountInString(input.Body) > 140 {
		SendJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Chirp is too long",
		})
		return
	}

	chirp, err := s.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		UserID: input.UserID,
		Body:   cleanBody(input.Body),
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	SendJSON(w, http.StatusOK, chirp)
}

var profane = []string{"kerfuffle", "sharbert", "fornax"}

func cleanBody(body string) string {
	words := strings.Fields(body)
	for i, word := range words {
		if slices.Contains(profane, strings.ToLower(word)) {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
