package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {

	requestChirpID := r.PathValue("chirpID")

	ctx := r.Context()

	chirpID, err := uuid.Parse(requestChirpID)
	if err != nil {
		log.Printf("failed to parse UUID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse UUID")
		return
	}

	dbChirp, err := cfg.db.GetChirp(ctx, chirpID)
	if err != nil {
		log.Printf("Error getting chirp from the database: %s", err)
		respondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}

	chirp := mapDbChirpStruct(dbChirp)

	respondWithJSON(w, http.StatusOK, chirp)
}
