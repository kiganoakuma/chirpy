package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kiganoakuma/chirpy/internal/auth"
	"github.com/kiganoakuma/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	// Only expect the body in the request
	type parameters struct {
		Body string `json:"body"`
	}

	// First, extract and validate the JWT token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing or invalid token", err)
		return
	}

	// Validate the token and get the user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtsecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: invalid token", err)
		return
	}

	// decode the request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validate the chirp content
	filteredBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	// Create the chirp in the database
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   filteredBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp", err)
		return
	}

	respose := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}

	// Return the created chirp
	respondWithJSON(w, http.StatusCreated, respose)
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	const replacmentChars = "****"
	filteredBody := replaceWords(body, profaneWords, replacmentChars)

	return filteredBody, nil
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
