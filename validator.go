package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chirpFields struct {
	Body         string `json:"body"`
	Valid        bool   `json:"valid"`
	Cleaned_Body string `json:"cleaned_body"`
	User_id      string `json:"user_id"`
}

func censorChirp(text string) string {
	const censorBar = "****"
	words := strings.Split(text, " ")
	for i, cheep := range words {
		switch strings.ToLower(cheep) {
		case "kerfuffle":
			fallthrough
		case "sharbert":
			fallthrough
		case "fornax":
			words[i] = censorBar
		default:
			continue
		}
	}
	cleanedChirp := strings.Join(words, " ")
	return cleanedChirp
}

func errorResponse(writer http.ResponseWriter, code int, text string) error {
	// too long error
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(code)
	if len(text) > 140 {
		writer.Write([]byte("Chirp is too long"))
		return fmt.Errorf("Chirp is too long")
	} else {
		writer.Write([]byte(text))
		return fmt.Errorf("%s", text)
	}
}

func jsonResponse(writer http.ResponseWriter, code int, form Cherp) {
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
