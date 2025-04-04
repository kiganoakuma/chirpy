package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kiganoakuma/chirpy/internal/auth"
	"github.com/kiganoakuma/chirpy/internal/database"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpID")

	chirpID, err := uuid.Parse(chirpIdString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp Id", err)
	}

	authHeader, respErr := auth.AuthHeaderHelper(w, r)
	if respErr != "" {
		respondWithError(w, http.StatusUnauthorized, respErr, nil)
		return
	}

	userID, err := auth.ValidateJWT(authHeader, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Unauthorized token", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if userID != chirp.UserID {
		respondWithError(w, http.StatusForbidden, "Not authorized to delete chirp", nil)
		return
	}

	err = cfg.db.DelChirpById(r.Context(), database.DelChirpByIdParams{
		ID:     chirpID,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete chirp", err)
	}
	respondWithJSON(w, http.StatusNoContent, nil)

}
