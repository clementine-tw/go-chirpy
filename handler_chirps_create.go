package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/clementine-tw/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleanedBody := replaceProfane(body)
	return cleanedBody, nil
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error decoding parameters",
			err,
		)
		return
	}

	if len(params.UserID) == 0 {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Missing user_id",
			err,
		)
		return
	}

	cleanedBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			err.Error(),
			err,
		)
		return
	}

	chirp, err := cfg.db.CreateChirp(
		r.Context(),
		database.CreateChirpParams{
			Body:   cleanedBody,
			UserID: params.UserID,
		},
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error creating chirp",
			err,
		)
		return
	}

	validBody := struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}{
		Body:   chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(
		w,
		http.StatusCreated,
		validBody,
	)
}
