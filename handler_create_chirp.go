package main

import (
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"

	"github.com/mmsacky/chirpy.git/internal/auth"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	const maxChirpLength = 140

	ctx := r.Context()

	type createChirpRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpRequest := createChirpRequest{}
	err := decoder.Decode(&chirpRequest)
	defer r.Body.Close()

	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	if utf8.RuneCountInString(chirpRequest.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header.")
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "token is invalid, tampered with, or expired.")
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   cleanChirp(chirpRequest.Body),
		UserID: userID,
	}

	dbChirp, err := cfg.db.CreateChirp(ctx, chirpParams)
	if err != nil {
		log.Printf("Error adding chirp to the database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	chirp := mapDbChirpStruct(dbChirp)

	respondWithJSON(w, http.StatusCreated, chirp)

}
