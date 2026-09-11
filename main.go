package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	serveM := http.NewServeMux()
	sever := http.Server{
		Addr:    ":8080",
		Handler: serveM,
	}

	serveM.Handle("/", http.FileServer(http.Dir(".")))

	err := sever.ListenAndServe()
	if err != nil {
		fmt.Println("ERROR server listening")
		os.Exit(1)
	}
}
