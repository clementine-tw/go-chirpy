package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error decoding json",
			err,
		)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Error creating user",
			err,
		)
		return
	}

	type userVal struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		userVal{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	)
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.db.ResetUsers(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error reseting users",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
}
