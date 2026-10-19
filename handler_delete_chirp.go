package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmsacky/chirpy.git/internal/auth"
)

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header.")
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "token is invalid, tampered with, or expired.")
		return
	}

	ctx := r.Context()

	requestChirpID := r.PathValue("chirpID")
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

	if dbChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You do not have permission to access this resource")
		return
	}

	err = cfg.db.DeleteChirp(ctx, dbChirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete chirp")
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
