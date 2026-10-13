package main

import (
	"net/http"

	"github.com/AirflowSteve/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "no token found", err)
		return
	}
	usrID, err := auth.ValidateJWT(bearer, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "incorrect token", err)
		return
	}
	chripID, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "no chirpid found", err)
		return
	}
	ch, err := cfg.db.GetChirpByID(r.Context(), chripID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "no chirp found", err)
		return
	}

	if usrID != ch.UserID {
		respondWithError(w, http.StatusForbidden, "forbidden to delete chirp", err)
		return
	}

	_, err = cfg.db.DeleteChirpByID(r.Context(), chripID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't delete the chirp", err)
	}

	w.WriteHeader(http.StatusNoContent)

}
