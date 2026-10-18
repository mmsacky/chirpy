package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbChirps, err := cfg.db.GetChirps(ctx)
	if err != nil {
		log.Printf("Error getting chirps from the database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {

		chirp := mapDbChirpStruct(dbChirp)
		chirps = append(chirps, chirp)

	}

	respondWithJSON(w, http.StatusOK, chirps)

}
