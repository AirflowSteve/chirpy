package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/AirflowSteve/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get the token", err)
		return
	}

	r_token, err := cfg.db.GetRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't get the token", err)
		return
	}
	if r_token.RevokedAt.Valid || (time.Now().Compare(r_token.ExpiresAt) != -1) {
		respondWithError(w, 401, "the refresh token is outdated", errors.New("refresh token ouf of date"))
		return
	}
	acc_token, err := auth.MakeJWT(r_token.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "could'n create a token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Token string `json:"token"`
	}{
		Token: acc_token,
	})

}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldn't get the token", err)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(r.Context(), bearer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revokee the token", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
