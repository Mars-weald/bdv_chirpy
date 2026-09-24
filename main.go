package main

import (
	"fmt"
	"net/http"
)

func main() {
	// server making
	serveM := http.NewServeMux()
	sever := http.Server{
		Addr:    ":8080",
		Handler: serveM,
	}

	apiConch := apiConfig{}

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
	serveM.HandleFunc("POST /api/validate_chirp", validateChirp)

	fmt.Printf("Serving files from . on port: 8080\n")
	sever.ListenAndServe()
}
