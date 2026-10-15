package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Mars-weald/bdv_chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	//database setup
	dbURL := os.Getenv("DB_URL")
	_, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("ERROR opening database")
		os.Exit(1)
	}
	plat := os.Getenv("PLATFORM")

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Printf("Error opening database: %s\n", err)
	}
	dbQueries := database.New(dbConn)

	// server making
	serveM := http.NewServeMux()
	sever := http.Server{
		Addr:    ":8080",
		Handler: serveM,
	}

	apiConch := apiConfig{
		qs:       dbQueries,
		is_admin: plat}

	serveM.Handle("/app/", apiConch.middlewareMetricsIncrement(http.FileServer(http.Dir("."))))
	// is working message
	pandle := func(writer http.ResponseWriter, req *http.Request) {
		req.Header.Set("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(200)
		writer.Write([]byte("OK"))
	}
	//register handlers
	serveM.HandleFunc("GET /api/healthz", pandle)
	serveM.HandleFunc("GET /admin/metrics", apiConch.serverHitsCount)
	serveM.HandleFunc("POST /admin/reset", apiConch.resetMetrics)
	serveM.HandleFunc("POST /api/chirps", apiConch.chirpEndpoint)
	serveM.HandleFunc("POST /api/users", apiConch.userHandler)

	fmt.Printf("Serving files from . on port: 8080\n")
	sever.ListenAndServe()
}
