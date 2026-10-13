package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/mmsacky/chirpy.git/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             database.Queries
	platform       string
}

type JSONerror struct {
	Error string `json:"error"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("Unable to connect to database")
	}
	dbQueries := database.New(db)

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Printf("Starting server at address localhost%s", server.Addr)

	apiCfg := &apiConfig{
		fileserverHits: atomic.Int32{},
		db:             *dbQueries,
		platform:       os.Getenv("PLATFORM"),
	}

	root := http.Dir(".")
	rootHandler := http.FileServer(root)

	mux.Handle("/app", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", rootHandler)))

	mux.HandleFunc("GET /api/healthz", readiness)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirp)
	mux.HandleFunc("POST /api/users", apiCfg.createUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirp)

	mux.HandleFunc("GET /admin/metrics", apiCfg.countReq)
	mux.HandleFunc("POST /admin/reset", apiCfg.reset)

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

func (cfg *apiConfig) reset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("403 Forbidden\n"))
		return
	}

	ctx := r.Context()
	cfg.db.ResetDB(ctx)
	cfg.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0 & all Database was reset\n"))

}

func (cfg *apiConfig) createChirp(w http.ResponseWriter, r *http.Request) {

	const maxChirpLength = 140

	ctx := r.Context()

	type createChirpRequest struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	chirpRequest := createChirpRequest{}
	err := decoder.Decode(&chirpRequest)
	defer r.Body.Close()

	if err != nil {
		log.Printf("Error decoding chirp: %w", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	if utf8.RuneCountInString(chirpRequest.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirpParams := database.CreateChirpParams{
		Body:   cleanChirp(chirpRequest.Body),
		UserID: chirpRequest.UserID,
	}

	dbChirp, err := cfg.db.CreateChirp(ctx, chirpParams)
	if err != nil {
		log.Printf("Error adding chirp to the database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {

	type createUserRequest struct {
		Email string `json:"email"`
	}
	ctx := r.Context()

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	userRequest := createUserRequest{}
	err := decoder.Decode(&userRequest)
	if err != nil {
		log.Printf("Error decoding user email address: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse JSON")
		return
	}

	dbUser, err := cfg.db.CreateUser(ctx, userRequest.Email)
	if err != nil {
		log.Printf("Error creating database user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dbChirps, err := cfg.db.GetChirps(ctx)
	if err != nil {
		log.Printf("Error getting chirps from the database: %s", err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get chirps")
		return
	}

	chirps := []Chirp{}

	for _, dbChirp := range dbChirps {
		chirp := Chirp{
			ID:        dbChirp.ID,
			Body:      dbChirp.Body,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
		}
		chirps = append(chirps, chirp)

	}
	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {

	requestChirpID := r.PathValue("chirpID")

	ctx := r.Context()

	chirpID, err := uuid.Parse(requestChirpID)
	if err != nil {
		log.Printf("failed to parse UUID: %s", err)
		respondWithError(w, http.StatusBadRequest, "Couldn't parse UUID")
		return
	}

	dbChirp, err := cfg.db.GetChirp(ctx, chirpID)
	if err != nil {
		log.Printf("Error getting chirp from the database: %s", err)
		respondWithError(w, http.StatusNotFound, "Chirp doesn't exist")
		return
	}

	chirp := Chirp{
		ID:        dbChirp.ID,
		Body:      dbChirp.Body,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
