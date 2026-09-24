package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type chirpFields struct {
	Body  string `json:"body"`
	Valid bool   `json:"valid"`
}

func errorResponse(writer http.ResponseWriter, code int, text string) error {
	// too long error
	if len(text) > 140 {
		writer.Header().Add("Content-Type", "text/plain; charset=utf-8")
		writer.WriteHeader(code)
		writer.Write([]byte("Chirp is too long"))
		return fmt.Errorf("Chirp is too long")
	}
	return nil
}

func jsonResponse(writer http.ResponseWriter, code int, form chirpFields) {
	datum, err := json.Marshal(form)
	if err != nil {
		writer.WriteHeader(500)
		fmt.Printf("ERROR marshaling json: %s\n", err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(datum)
}

func validateChirp(writer http.ResponseWriter, req *http.Request) {
	decoderPin := json.NewDecoder(req.Body)
	paramaters := chirpFields{}
	err := decoderPin.Decode(&paramaters)
	if err != nil {
		fmt.Printf("ERROR deconding json: %s\n", err)
		writer.WriteHeader(500)
		return
	}

	err = errorResponse(writer, 400, paramaters.Body)
	if err != nil {
		return
	}

	paramaters.Valid = true
	jsonResponse(writer, 200, paramaters)
}
