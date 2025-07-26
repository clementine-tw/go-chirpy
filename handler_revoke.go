package main

import (
	"net/http"

	"github.com/clementine-tw/go-chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshTokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't get refresh token",
			err,
		)
		return
	}

	_, err = cfg.db.RevokeRefreshToken(
		r.Context(),
		refreshTokenString,
	)
	if err != nil {
		respondWithError(
			w,
			http.StatusUnauthorized,
			"Coudn't find refresh token",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
