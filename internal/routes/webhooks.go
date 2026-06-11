package routes

import (
	"net/http"
	"os"

	"github.com/faymndev/chirpy/internal/auth"
	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
	"github.com/google/uuid"
)

func UseWebhooks(mux *http.ServeMux, state *middleware.State) {
	wh := webhooks{state: state}
	mux.HandleFunc("POST /api/polka/webhooks", wh.handlePolkaEvent)
}

type webhooks struct {
	state *middleware.State
}

func (wh *webhooks) handlePolkaEvent(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiToken(r.Header)
	if err != nil || apiKey != os.Getenv("POLKA_API_KEY") {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid API Key",
		})
		return
	}

	type Input struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	input, err := DecodeBody[Input](r)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Something went wrong",
		})
		return
	} else if input.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userID, err := uuid.Parse(input.Data.UserID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, err = wh.state.Db.GetUserById(r.Context(), userID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	_, err = wh.state.Db.UpgradeUser(r.Context(), database.UpgradeUserParams{
		ID:          userID,
		IsChirpyRed: true,
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
