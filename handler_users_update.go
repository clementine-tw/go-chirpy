package main

import (
	"encoding/json"
	"net/http"

	"github.com/clementine-tw/go-chirpy/internal/auth"
	"github.com/clementine-tw/go-chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
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
			http.StatusUnauthorized,
			"Couldn't validate JWT",
			err,
		)
		return
	}

	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameter{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't decode body",
			err,
		)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't hash password",
			err,
		)
		return
	}

	updatedUser, err := cfg.db.UpdateUserEmailAndPassword(
		r.Context(),
		database.UpdateUserEmailAndPasswordParams{
			Email:          params.Email,
			HashedPassword: hashedPassword,
			ID:             userID,
		},
	)

	newJWT, err := auth.MakeJWT(updatedUser.ID, cfg.secret, defaultJWTExpiration)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't make JWT",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		User{
			ID:          updatedUser.ID,
			CreatedAt:   updatedUser.CreatedAt,
			UpdatedAt:   updatedUser.UpdatedAt,
			Email:       updatedUser.Email,
			Token:       newJWT,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	)

}
