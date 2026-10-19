package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/mmsacky/chirpy.git/internal/auth"
)

func (cfg *apiConfig) getRefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header.")
		return
	}

	dbRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusUnauthorized, "token is invalid, tampered with, or expired.")
		} else {
			log.Printf("Error getting user from the database: %v", err)
			respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		}
		return
	}

	if dbRefreshToken.RevokedAt.Valid || !dbRefreshToken.ExpiresAt.After(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "token is invalid, tampered with, or expired.")
		return
	}

	accessToken, err := auth.MakeJWT(dbRefreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Secure session could not be established. Please try again.")
		return
	}

	token := mapAccessTokenStruct(accessToken)

	respondWithJSON(w, http.StatusOK, token)
}
