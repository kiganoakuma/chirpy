package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/kiganoakuma/chirpy/internal/auth"
	"github.com/kiganoakuma/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameter", err)
		return
	}

	// Get user from database
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Verify password before creating token
	err = auth.CheckPassword(user.HashedPassword, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Determine expiration time (default: 1 hour)
	jwtExpTime := time.Hour

	// Create JWT token
	tokenString, err := auth.MakeJWT(
		user.ID,
		cfg.jwtsecret,
		jwtExpTime,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Authentication token", err)
		return
	}
	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't generate refresh token", err)
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 60)

	// Create refresh token
	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshTokenString,
		UserID:    user.ID,
		ExpiresAt: refreshTokenExp,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't insert token into database", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       params.Email,
			IsChirpyRed: user.IsChirpyRed.Bool,
		},
		Token:        tokenString,
		RefreshToken: refreshToken.Token,
	})
}
