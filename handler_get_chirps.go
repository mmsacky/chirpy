package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/mmsacky/chirpy.git/internal/database"
)

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	var err error
	var dbChirps []database.Chirp
	ctx := r.Context()
	authorIDString := r.URL.Query().Get("author_id")
	sortBy := r.URL.Query().Get("sort")

	if authorIDString != "" {
		authorID, parseErr := uuid.Parse(authorIDString)
		if parseErr != nil {
			log.Printf("failed to parse UUID: %s", err)
			respondWithError(w, http.StatusBadRequest, "Couldn't parse UUID")
			return
		}
		dbChirps, err = cfg.db.GetChirpsByUserID(ctx, authorID)
	} else {
		dbChirps, err = cfg.db.GetChirps(ctx)
	}
	if err != nil {
		log.Printf("Error getting chirps from the database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := make([]Chirp, 0, len(dbChirps))

	for _, dbChirp := range dbChirps {
		chirps = append(chirps, mapDbChirpStruct(dbChirp))
	}

	if sortBy == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})

	}

	respondWithJSON(w, http.StatusOK, chirps)

}
