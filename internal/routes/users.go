package routes

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/faymndev/chirpy/internal/auth"
	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
)

func UseUsers(mux *http.ServeMux, state *middleware.State) {
	mux.Handle("POST /api/login", state.Middleware(handleLogin))
	mux.Handle("POST /api/refresh", state.Middleware(handleRefresh))
	mux.Handle("POST /api/revoke", state.Middleware(handleRevoke))
	mux.Handle("POST /api/users", state.Middleware(handleCreateUser))
}

func handleLogin(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	input := Input{}
	if err := decoder.Decode(&input); err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	user, err := s.Db.GetUser(r.Context(), input.Email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	match, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	} else if !match {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid username or password combination",
		})
		return
	}

	expiresIn, _ := time.ParseDuration("1h")
	token, err := auth.MakeJWT(user.ID, expiresIn)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	// do we already have a refresh token?
	refreshToken, err := s.Db.GetUserRefreshToken(r.Context(), user.ID)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	var newRefreshToken string
	if auth.IsValidRefreshToken(refreshToken) {
		newRefreshToken = refreshToken.Token
	} else {
		newRefreshToken = auth.MakeRefreshToken()
		err = s.Db.SetRefreshToken(r.Context(), database.SetRefreshTokenParams{
			Token:  newRefreshToken,
			UserID: user.ID,
		})
		if err != nil {
			SendJSON(w, http.StatusInternalServerError, map[string]any{
				"error": "Something went wrong",
			})
		}
	}

	type UserWithToken struct {
		database.User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	SendJSON(w, http.StatusOK, UserWithToken{User: user, Token: token, RefreshToken: newRefreshToken})
}

func handleRefresh(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Refresh token not provided",
		})
		return
	}

	refreshToken, err := s.Db.GetRefreshToken(r.Context(), token)
	if err != nil || !auth.IsValidRefreshToken(refreshToken) {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid refresh token",
		})
		return
	}

	newToken := auth.MakeRefreshToken()
	err = s.Db.SetRefreshToken(r.Context(), database.SetRefreshTokenParams{
		Token:  newToken,
		UserID: refreshToken.UserID,
	})
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Failed to refresh token",
		})
		return
	}

	SendJSON(w, http.StatusOK, map[string]any{
		"token": newToken,
	})
}

func handleRevoke(w http.ResponseWriter, r *http.Request, s *middleware.State) {

}

func handleCreateUser(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

	hashedPassword, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to hash password",
		})
	}

	user, err := s.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
	}

	SendJSON(w, http.StatusCreated, user)
}
