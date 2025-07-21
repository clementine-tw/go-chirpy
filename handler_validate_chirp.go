package main

import (
	"encoding/json"
	"net/http"
)

const chirpCharsLimitLen = 140

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
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

	if len(params.Body) > chirpCharsLimitLen {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Chirp is too long",
			err,
		)
		return
	}

	validBody := struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	}

	respondWithJSON(
		w,
		http.StatusOK,
		validBody,
	)
}
