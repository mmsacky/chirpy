package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mmsacky/chirpy.git/internal/auth"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	type createUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	userRequest := createUserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Printf("Error decoding user information: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	hash, err := auth.HashPassword(userRequest.Password)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Unexpected error")
		return
	}

	userParams := database.CreateUserParams{
		Email:          userRequest.Email,
		HashedPassword: hash,
	}

	dbUser, err := cfg.db.CreateUser(ctx, userParams)
	if err != nil {
		log.Printf("Error creating database user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	user := mapDbUserStruct(dbUser, "", "")

	respondWithJSON(w, http.StatusCreated, user)

}
