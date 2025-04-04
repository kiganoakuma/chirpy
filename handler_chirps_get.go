package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kiganoakuma/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorIdString := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error

	if authorIdString != "" {
		authorID, parseErr := uuid.Parse(authorIdString)
		if parseErr != nil {
			respondWithError(w, http.StatusBadRequest, "invalid id format", parseErr)
			return
		}
		chirps, err = cfg.db.GetChirpByAuthor(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could not retrieve chirps", err)
		return
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
