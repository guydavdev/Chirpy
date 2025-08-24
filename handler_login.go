package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/guydavdev/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to get user", err)
		return
	}

	hasAuthError := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if hasAuthError != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", hasAuthError)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
