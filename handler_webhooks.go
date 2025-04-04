package main

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/kiganoakuma/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type paramaters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No api key found", err)
		return
	}

	if cfg.polkakey != apiKey {
		respondWithError(w, 401, "Unauthorized api key", nil)
	}

	params, err := DecodeJSONBody[paramaters](r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode request", err)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid id paramater", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithError(w, 204, "", nil)
		return
	}

	err = cfg.db.UpgradeToRedByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	respondWithJSON(w, 204, "User upgraded to red")

}
