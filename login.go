package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/AirflowSteve/chirpy/internal/auth"
	"github.com/AirflowSteve/chirpy/internal/database"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	usr := parameters{}

	err := decoder.Decode(&usr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode the request", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), usr.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(usr.Password, user.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	tkn, err := auth.MakeJWT(user.ID, cfg.secret, time.Duration(3600*time.Second))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could'n create a token", err)
		return
	}
	refresh_token := auth.MakeRefreshToekn()

	r_tkn, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(1440 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create the refresh token", err)
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tkn,
		RefreshToken: r_tkn.Token,
	})
}
