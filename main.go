package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"github.com/AirflowSteve/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func handlerHealthz(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
<html>

<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, 403, "You're not allowed to do that", errors.New("not authorized to do things"))
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset the database", err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}

	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)

}

func main() {

	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	plfrm := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Printf("there was an error", err)
	}

	DBQueries := database.New(db)

	handler := http.NewServeMux()
	server := http.Server{
		Handler: handler,
		Addr:    ":8080",
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             DBQueries,
		platform:       plfrm,
	}

	handler.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	handler.HandleFunc("GET /api/healthz", handlerHealthz)
	handler.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	handler.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	handler.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	handler.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	handler.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	server.ListenAndServe()
}
