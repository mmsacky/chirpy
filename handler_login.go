package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/mmsacky/chirpy.git/internal/auth"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defaultExpirationSeconds := 3600

	type createUserRequest struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	userRequest := createUserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Printf("Error decoding user information: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(ctx, userRequest.Email)
	if err != nil {
		log.Printf("Error getting user from the database: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	match, err := auth.CheckPasswordHash(userRequest.Password, dbUser.HashedPassword)
	if err != nil {
		log.Printf("Error comparing password")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !match {
		log.Printf("Error password doesn't match")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	if userRequest.ExpiresInSeconds == 0 || userRequest.ExpiresInSeconds > defaultExpirationSeconds {
		userRequest.ExpiresInSeconds = defaultExpirationSeconds
	}

	expiresIn := time.Duration(userRequest.ExpiresInSeconds) * time.Second

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, expiresIn)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Secure session could not be established. Please try again.")
		return
	}
	dbUser.Token = token
	user := mapDbUserStruct(dbUser)

	respondWithJSON(w, http.StatusOK, user)

}
