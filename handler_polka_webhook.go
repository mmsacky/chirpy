package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/mmsacky/chirpy.git/internal/auth"
)

func (cfg *apiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error encountered: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Missing Authorization header.")
		return
	}

	if apiKey != cfg.polkaAPIKey {
		respondWithError(w, http.StatusUnauthorized, "Incorrect API key provided.")
		return
	}

	type chirpyRedRequest struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	chirpyRedReq := chirpyRedRequest{}
	err = decoder.Decode(&chirpyRedReq)
	if err != nil {
		log.Printf("Error decoding chirpy Red information: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	if chirpyRedReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	affectedRows, err := cfg.db.JoinChirpyRed(r.Context(), chirpyRedReq.Data.UserID)
	if err != nil {
		log.Printf("Error updating user information: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if affectedRows == 0 {
		respondWithError(w, http.StatusNotFound, "Resource not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
