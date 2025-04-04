package main

import (
	"net/http"
	"time"

	"github.com/kiganoakuma/chirpy/internal/auth"
)

type response struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {

	authHeader, respErr := auth.AuthHeaderHelper(w, r)
	if respErr != "" {
		respondWithError(w, http.StatusUnauthorized, respErr, nil)
		return
	}

	refreshToken, err := cfg.db.GetRefreshToken(r.Context(), authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}
	user, err := cfg.db.GetUserForRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(
		user.ID,
		cfg.jwtsecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create token", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: newAccessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader, respErr := auth.AuthHeaderHelper(w, r)
	if respErr != "" {
		respondWithError(w, http.StatusUnauthorized, respErr, nil)
		return
	}

	err := cfg.db.RevokeRefreshToken(r.Context(), authHeader)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
