package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

// server hits stuff
func (conf *apiConfig) middlewareMetricsIncrement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		conf.fileServerHits.Add(1)
		next.ServeHTTP(writer, req)
	})
}

func (conf *apiConfig) serverHitsCount(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	hits := conf.fileServerHits.Load()
	x := fmt.Sprintf("Hits: %v", hits)
	writer.Write([]byte(x))

}

func (conf *apiConfig) resetMetrics(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(200)
	conf.fileServerHits.Store(0)
	writer.Write([]byte("Server hits reset"))
}
