package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/AirflowSteve/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	apiKey         string
}

func handlerHealthz(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/plain; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(http.StatusText(http.StatusOK)))
}

func main() {

	godotenv.Load(".env")
	dbUrl := os.Getenv("DB_URL")
	plfrm := os.Getenv("PLATFORM")
	scrt := os.Getenv("SECRET")
	apiKEY := os.Getenv("PAYMENTAPI")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Printf("there was an error", err)
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
		secret:         scrt,
		apiKey:         apiKEY,
	}

	handler.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	handler.HandleFunc("GET /api/healthz", handlerHealthz)
	handler.HandleFunc("GET /admin/metrics", apiCfg.handlerHits)
	handler.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	handler.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	handler.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	handler.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	handler.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpByID)
	handler.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	handler.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	handler.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	handler.HandleFunc("PUT /api/users", apiCfg.handlerChangeCredentials)
	handler.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)
	handler.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerUpgradeUser)
	server.ListenAndServe()
}
