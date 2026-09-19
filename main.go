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
	serveM.HandleFunc("GET /healthz", pandle)
	serveM.HandleFunc("GET /metrics", apiConch.serverHitsCount)
	serveM.HandleFunc("POST /reset", apiConch.resetMetrics)

	fmt.Printf("Serving files from . on port: 8080\n")
	sever.ListenAndServe()
}
