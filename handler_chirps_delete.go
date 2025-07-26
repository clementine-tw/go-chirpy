package main

import (
	"net/http"

	"github.com/clementine-tw/go-chirpy/internal/auth"
	"github.com/clementine-tw/go-chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			"Couldn't validate JWT",
			err,
		)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't parse chirp id",
			err,
		)
		return
	}

	_, err = cfg.db.GetChirp(
		r.Context(),
		chirpID,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	_, err = cfg.db.DeleteChirp(
		r.Context(),
		database.DeleteChirpParams{
			ID:     chirpID,
			UserID: userID,
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusForbidden,
			"Not the owner of chirp",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
