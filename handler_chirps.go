package main

import (
	"encoding/json"
	"net/http"
)

func (cfg apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	type parameters struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	type returnVals struct {
		Cleaned string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	if params.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "UserId must be set", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{
		Cleaned: cleanBody(params.Body),
	})
}
