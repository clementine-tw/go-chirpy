package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/clementine-tw/go-chirpy/internal/auth"
)

const defaultExpirationDuration = 1 * time.Hour

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameter struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error decoding parameter",
			err,
		)
		return
	}

	if len(params.Password) == 0 || len(params.Email) == 0 {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Missing password or email",
			nil,
		)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Incorrect email or password",
			err,
		)
		return
	}

	expirationDuration := defaultExpirationDuration
	if params.ExpiresInSeconds > 0 {
		expirationDuration = time.Duration(params.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationDuration)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error making JWT",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Token:     token,
		},
	)
}
