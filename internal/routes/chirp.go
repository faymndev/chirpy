package routes

import (
	"encoding/json"
	"net/http"
	"slices"
	"strings"
	"unicode/utf8"
)

func UseChirp(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/validate_chirp", handleValidateChirp)
}

func handleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Body string `json:"body"`
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

	SendJSON(w, http.StatusOK, map[string]any{
		"cleaned_body": cleanBody(input.Body),
	})
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
