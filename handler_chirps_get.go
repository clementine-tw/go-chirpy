package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {

	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(
			w,
			http.StatusBadRequest,
			"Error parse chirp id",
			err,
		)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(
			w,
			http.StatusNotFound,
			"Chirp not found",
			err,
		)
		return
	}

	respondWithJSON(
		w,
		http.StatusOK,
		Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	)
}

func (cfg *apiConfig) handlerChirpsList(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(
			w,
			http.StatusInternalServerError,
			"Error quering chirps",
			err,
		)
		return
	}

	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			respondWithError(
				w,
				http.StatusBadRequest,
				"Invalid user id",
				err,
			)
			return
		}
	}

	val := []Chirp{}
	for _, chirp := range chirps {
		if authorID != uuid.Nil && authorID != chirp.UserID {
			continue
		}
		val = append(val, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	sortType := r.URL.Query().Get("sort")
	if sortType == "desc" {
		sort.Slice(val, func(i, j int) bool {
			return val[i].CreatedAt.After(val[j].CreatedAt)
		})
	}

	respondWithJSON(
		w,
		http.StatusOK,
		val,
	)
}
