package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/mmsacky/chirpy.git/internal/auth"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now()

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header.")
		return
	}

	revokeRefreshTokenParams := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: now,
		Token:     refreshToken,
	}

	err = cfg.db.RevokeRefreshToken(ctx, revokeRefreshTokenParams)
	if err != nil {
		log.Printf("Error revoking refresh token: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
