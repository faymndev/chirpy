package routes

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/faymndev/chirpy/internal/auth"
	"github.com/faymndev/chirpy/internal/database"
	"github.com/faymndev/chirpy/internal/middleware"
)

func UseUsers(mux *http.ServeMux, state *middleware.State) {
	mux.Handle("POST /api/login", state.Middleware(handleLogin))
	mux.Handle("POST /api/refresh", state.Middleware(handleRefresh))
	mux.Handle("POST /api/revoke", state.Middleware(handleRevoke))
	mux.Handle("POST /api/users", state.Middleware(handleCreateUser))
	mux.Handle("PUT /api/users", state.Middleware(handleUpdateUser))
}

func handleLogin(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input, err := DecodeBody[Input](r)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to decode request body",
		})
		return
	}

	user, err := s.Db.GetUser(r.Context(), input.Email)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to get user by email",
		})
		return
	}

	err = auth.ComparePassword(input.Password, user.Password)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid username or password combination",
		})
		return
	}

	token, err := auth.MakeJWT(user.ID, time.Hour)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to create JWT token",
		})
		return
	}

	refreshToken, err := s.Db.GetUserRefreshToken(r.Context(), user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		refreshToken, err = s.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
			Token:     auth.MakeRefreshToken(),
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60), // expire in 60 days
		})
	}
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to retrieve refresh token",
		})
		return
	}

	// reuse refresh token if valid, otherwise, create a new one
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
			return
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

	newToken, err := auth.MakeJWT(refreshToken.UserID, time.Hour)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to create JWT token",
		})
		return
	}

	SendJSON(w, http.StatusOK, map[string]any{
		"token": newToken,
	})
}

func handleRevoke(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Refresh token not provided",
		})
		return
	}

	if _, err = s.Db.GetRefreshToken(r.Context(), token); err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid refresh token",
		})
		return
	}

	if err = s.Db.RevokeRefreshToken(r.Context(), token); err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to revoke refresh token",
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input, err := DecodeBody[Input](r)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Failed to hash password",
		})
		return
	}

	user, err := s.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	SendJSON(w, http.StatusCreated, user)
}

func handleUpdateUser(w http.ResponseWriter, r *http.Request, s *middleware.State) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid token",
		})
		return
	}

	userId, err := auth.VerifyJWT(token)
	if err != nil {
		SendJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "Invalid token",
		})
		return
	}

	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	input, err := DecodeBody[Input](r)
	if err != nil {
		SendJSON(w, http.StatusBadRequest, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	user, err := s.Db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:       userId,
		Email:    input.Email,
		Password: hashedPassword,
	})
	if err != nil {
		SendJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "Something went wrong",
		})
		return
	}

	SendJSON(w, http.StatusOK, user)
}
