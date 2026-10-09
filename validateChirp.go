package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidate(w http.ResponseWriter, r *http.Request) {

	type info struct {
		Body string `json:"body"`
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	decoder := json.NewDecoder(r.Body)

	data := info{}

	err := decoder.Decode(&data)
	if err != nil {
		respondWithError(w, 400, "Something went wrong")
		return
	}

	words := strings.Split(data.Body, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = "****"
		}

	}
	data.Body = strings.Join(words, " ")

	if len(data.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: data.Body,
	})

}
