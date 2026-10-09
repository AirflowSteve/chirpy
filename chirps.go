package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AirflowSteve/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirp struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Body       string    `json:"body"`
	UserID     uuid.UUID `json:"user_id"`
}

func validateChirp(c string) (string, error) {
	const maxChirpLength = 140

	if len(c) > maxChirpLength {
		return "", fmt.Errorf("chirp is too long! %v\n", len(c))
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	words := strings.Split(c, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}

	}
	result := strings.Join(words, " ")

	return result, nil
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", err)
		return
	}
	c, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "chirp wasn't validated", err)
		return
	}
	result_ch, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   c,
		UserID: params.UserID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a chirp", err)
		return
	}

	respondWithJSON(w, 201, chirp{
		ID:         result_ch.ID,
		Created_at: result_ch.CreatedAt,
		Updated_at: result_ch.UpdatedAt,
		Body:       result_ch.Body,
		UserID:     result_ch.UserID,
	})

}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirps", err)
	}

	result := []chirp{}
	for _, chirp_ := range chirps {
		result = append(result, chirp{
			ID:         chirp_.ID,
			Created_at: chirp_.CreatedAt,
			Updated_at: chirp_.UpdatedAt,
			Body:       chirp_.Body,
			UserID:     chirp_.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, result)
}
