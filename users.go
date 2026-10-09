package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type rBody struct {
		Email string `json:"email"`
	}

	result := rBody{}

	err := decoder.Decode(&result)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}
	usr, err := cfg.db.CreateUser(r.Context(), result.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create the user", err)
		return
	}

	user := User{
		ID:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}

	respondWithJSON(w, 201, user)

}
