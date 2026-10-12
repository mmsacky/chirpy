package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

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
