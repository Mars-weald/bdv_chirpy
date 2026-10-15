package main

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/Mars-weald/bdv_chirpy/internal/database"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	qs             *database.Queries
	is_admin       string
}

// server hits stuff
func (conf *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		conf.fileServerHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}

func (conf *apiConfig) serverHitsCount(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(200)
	hits := conf.fileServerHits.Load()
	x := fmt.Sprintf("<html> <body> <h1>Welcome, Chirpy Admin</h1> <p>Chirpy has been visited %d times!</p> </body </html>", hits)
	writer.Write([]byte(x))

}

func (conf *apiConfig) resetMetrics(writer http.ResponseWriter, req *http.Request) {
	if conf.is_admin != "dev" {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(403)
		writer.Write([]byte("FORBIDDEN COMMAND"))
	} else {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(200)
		conf.fileServerHits.Store(0)
		conf.qs.ResetUsers(req.Context())
		writer.Write([]byte("Server hits, users table reset"))
	}
}
