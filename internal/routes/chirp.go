package routes

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/faymndev/chirpy/internal/auth"
	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
	"github.com/google/uuid"
)

func UseChirp(mux *http.ServeMux, state *middleware.State) {
	mux.Handle("GET /api/chirps", state.Middleware(handleGetChirps))
	mux.Handle("GET /api/chirps/{chirpID}", state.Middleware(handleGetChirp))
	mux.Handle("POST /api/chirps", state.Middleware(handleChirp))
}

func handleGetChirp(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		SendJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Invalid chirp ID",
		})
		return
	}

	chirp, err := s.Db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}

		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	SendJSON(w, http.StatusOK, chirp)
}

func handleGetChirps(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	chirps, err := s.Db.GetChirps(r.Context())
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
	}

	SendJSON(w, http.StatusOK, chirps)
}

func handleChirp(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "JWT token not provided",
		})
		return
	}

	userID, err := auth.VerifyJWT(token)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid JWT token",
		})
		return
	}

	type Input struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	input := Input{}
	err = decoder.Decode(&input)
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
		UserID: userID,
		Body:   cleanBody(input.Body),
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	SendJSON(w, http.StatusCreated, chirp)
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
