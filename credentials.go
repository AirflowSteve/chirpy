package main

import (
	"encoding/json"
	"net/http"

	"github.com/AirflowSteve/chirpy/internal/auth"
	"github.com/AirflowSteve/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChangeCredentials(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token found", err)
	}

	userID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid access token", err)
	}
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type user_parameters struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "blya", err)
	}

	hashed_password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update the password", err)
	}

	usr, err := cfg.db.UpdateCreds(r.Context(), database.UpdateCredsParams{
		Email:          params.Email,
		HashedPassword: hashed_password,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update the password", err)
	}

	respondWithJSON(w, 200, user_parameters{
		ID:          usr.ID,
		CreatedAt:   usr.CreatedAt,
		Email:       usr.Email,
		Token:       bearer,
		IsChirpyRed: usr.IsChirpyRed,
		UpdatedAt:   usr.UpdatedAt,
	})

}
