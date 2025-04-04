package main

import (
	"encoding/json"
	"net/http"

	"github.com/kiganoakuma/chirpy/internal/auth"
	"github.com/kiganoakuma/chirpy/internal/database"
)

type paramaters struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {

	authHeader, respErr := auth.AuthHeaderHelper(w, r)
	if respErr != "" {
		respondWithError(w, http.StatusUnauthorized, respErr, nil)
		return
	}

	userID, err := auth.ValidateJWT(authHeader, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := paramaters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't not parse response body", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	user, err := cfg.db.PutUserCredentials(r.Context(), database.PutUserCredentialsParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}
