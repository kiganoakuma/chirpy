package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kiganoakuma/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	const replacmentChars = "****"

	filteredBody := replaceWords(params.Body, profaneWords, replacmentChars)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   filteredBody,
		UserID: params.UserId,
	})
	if err != nil {
		log.Fatalf("could not create chirp: %s", err)
	}

	respondWithJSON(w, http.StatusCreated, response{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    params.UserId,
		},
	})
}

func replaceWords(input string, profanity []string, replacement string) string {
	escapedTargets := make([]string, len(profanity))
	for i, word := range profanity {
		escapedTargets[i] = regexp.QuoteMeta(word)
	}

	pattern := `(?i)\b(` + strings.Join(escapedTargets, "|") + `)\b`
	re := regexp.MustCompile(pattern)
	result := re.ReplaceAllString(input, replacement)

	return result
}
