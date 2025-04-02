package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Fatalf("could not retrieve chirps: %s", err)
	}

	chirpsSlice := make([]Chirp, 0, len(chirps))

	for _, chirp := range chirps {
		chirpsSlice = append(chirpsSlice, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		})

	}
	respondWithJSON(w, http.StatusOK, chirpsSlice)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Chirp
	}

	chirpIdString := r.PathValue("chirpID")

	chirpId, err := uuid.Parse(chirpIdString)
	if err != nil {
		log.Fatalf("could not parse uuid string: %s", err)
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		respondWithError(w, 404, "Chirp does not exist", err)
	}

	respondWithJSON(w, http.StatusOK, response{
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		},
	})
}
