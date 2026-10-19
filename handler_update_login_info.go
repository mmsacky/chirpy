package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mmsacky/chirpy.git/internal/auth"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) updateLoginInfo(w http.ResponseWriter, r *http.Request) {

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

	type updateRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	userUpdateRequest := updateRequest{}
	err = decoder.Decode(&userUpdateRequest)
	if err != nil {
		log.Printf("Error decoding user information: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	hashedPassword, err := auth.HashPassword(userUpdateRequest.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	loginInfoParams := database.UpdateLoginInfoParams{
		Email:          userUpdateRequest.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	}

	dbUser, err := cfg.db.UpdateLoginInfo(ctx, loginInfoParams)
	if err != nil {
		log.Printf("Error updating user information: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	user := mapDbUserStruct(dbUser, "", "")

	respondWithJSON(w, http.StatusOK, user)
}
