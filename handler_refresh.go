package main

import (
	"net/http"

	"github.com/clementine-tw/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't get bearer token",
			err,
		)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(
		r.Context(),
		refreshTokenString,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Couldn't find user",
			err,
		)
		return
	}

	newJWT, err := auth.MakeJWT(user.ID, cfg.secret, defaultJWTExpiration)
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Couldn't make JWT",
			err,
		)
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: newJWT,
	}

	respondWithJSON(
		w,
		http.StatusOK,
		resp,
	)
}
