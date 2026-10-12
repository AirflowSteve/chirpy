package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AirflowSteve/chirpy/internal/auth"
	"github.com/AirflowSteve/chirpy/internal/database"
	"github.com/google/uuid"
)

type RBody struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	result := RBody{}

	err := decoder.Decode(&result)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode the request", err)
		return
	}

	hashed_pswrd, err := auth.HashPassword(result.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error setting the password", err)
		return
	}

	usr, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          result.Email,
		HashedPassword: hashed_pswrd,
	})
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
