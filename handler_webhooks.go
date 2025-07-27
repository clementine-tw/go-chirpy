package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Couldn't decode body",
			err,
		)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
	}

	_, err = cfg.db.UpgradeUserChirpyRed(
		r.Context(),
		params.Data.UserID,
	)

	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"User not found",
			err,
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
