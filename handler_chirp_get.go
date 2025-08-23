package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	if chirpID == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID", nil)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirpID", nil)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "Chirp not found", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get Chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
