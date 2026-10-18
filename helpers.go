package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mmsacky/chirpy.git/internal/database"
)

type JSONerror struct {
	Error string `json:"error"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	resp := JSONerror{
		Error: msg,
	}

	respData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respData)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	respData, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respData)
}

func cleanChirp(chirp string) string {

	badWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Fields(chirp)

	for i := range words {
		if slices.Contains(badWords, strings.ToLower(words[i])) {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}

func mapDbChirpStruct(dbChirp database.Chirp) Chirp {

	return Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
	}
}

func mapDbUserStruct(dbUser database.User) User {

	return User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     dbUser.Token,
	}
}
