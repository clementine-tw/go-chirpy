package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/clementine-tw/go-chirpy/internal/auth"
	"github.com/clementine-tw/go-chirpy/internal/database"
)

const (
	defaultJWTExpiration = 1 * time.Hour

	defaultRefreshTokenExpiration = 60 * 24 * time.Hour
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameter struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, defaultJWTExpiration)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error making JWT",
			err,
		)
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't make refresh token",
			err,
		)
		return
	}

	refreshTokenRecord, err := cfg.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshTokenString,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserID:    user.ID,
			ExpiresAt: time.Now().Add(defaultRefreshTokenExpiration),
			RevokedAt: sql.NullTime{
				Valid: false,
			},
		},
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't create refresh token",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:           user.ID,
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
			Email:        user.Email,
			Token:        token,
			RefreshToken: refreshTokenRecord.Token,
			IsChirpyRed:  user.IsChirpyRed,
		},
	)
}
