package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mmsacky/chirpy.git/internal/auth"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	now := time.Now()

	type createLoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	loginRequest := createLoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		log.Printf("Error decoding user information: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(ctx, loginRequest.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		} else {
			log.Printf("Error getting user from the database: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	match, err := auth.CheckPasswordHash(loginRequest.Password, dbUser.HashedPassword)
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

	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Secure session could not be established. Please try again.")
		return
	}

	refreshToken := auth.MakeRefreshToken()

	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    dbUser.ID,
		ExpiresAt: now.UTC().Add(60 * 24 * time.Hour),
		RevokedAt: sql.NullTime{},
	}

	_, err = cfg.db.CreateRefreshToken(ctx, refreshTokenParams)
	if err != nil {
		log.Printf("Error creating database refresh Token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Secure session could not be established. Please try again.")
		return
	}

	user := mapDbUserStruct(dbUser, accessToken, refreshToken)

	respondWithJSON(w, http.StatusOK, user)

}
