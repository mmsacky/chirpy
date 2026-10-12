package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

type JSONerror struct {
	Error string `json:"error"`
}

type validChirp struct {
	CleanedBody string `json:"cleaned_body"`
}

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
	}

	root := http.Dir(".")
	rootHandler := http.FileServer(root)

	mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", rootHandler)))

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.countReq)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetCountReq)

	log.Fatal(server.ListenAndServe())
}

func readiness(w http.ResponseWriter, _ *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}

func (cfg *apiConfig) countReq(w http.ResponseWriter, _ *http.Request) {

	count := cfg.fileserverHits.Load()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", count)))
}

func (cfg *apiConfig) resetCountReq(w http.ResponseWriter, _ *http.Request) {

	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte("Hits reset to 0\n"))
}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	type userPostRequest struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := userPostRequest{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedChirp := cleanChirp(chirp.Body)

	validChirp := validChirp{
		CleanedBody: cleanedChirp,
	}

	respondWithJSON(w, http.StatusOK, validChirp)

}

// Middleware

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
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
