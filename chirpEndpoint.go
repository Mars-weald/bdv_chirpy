package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Mars-weald/bdv_chirpy/internal/database"
	"github.com/google/uuid"
)

type Cherp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	User_id   uuid.UUID `json:"user_id"`
}

func (conf apiConfig) chirpEndpoint(writer http.ResponseWriter, req *http.Request) {
	decoderPin := json.NewDecoder(req.Body)
	paramaters := chirpFields{}
	err := decoderPin.Decode(&paramaters)
	if err != nil {
		fmt.Printf("ERROR deconding json: %s\n", err)
		writer.WriteHeader(500)
		return
	}

	if len(paramaters.Body) > 140 {
		errorResponse(writer, 400, paramaters.Body)
	}

	paramaters.Cleaned_Body = censorChirp(paramaters.Body)

	paramaters.Valid = true

	//Add to database
	oozer, err := uuid.Parse(paramaters.User_id)
	if err != nil {
		errorResponse(writer, 400, "ERROR parsing uuid")
		return
	}
	addChirpParams := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      paramaters.Cleaned_Body,
		UserID:    oozer,
	}
	chirp, err := conf.qs.CreateChirp(context.Background(), addChirpParams)
	if err != nil {
		errorResponse(writer, 400, "ERROR creating chirp entry")
		return
	}

	//tag fixing
	fixedtags := Cherp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	jsonResponse(writer, 201, fixedtags)
}

func (conf apiConfig) getChirpsEndpoint(writer http.ResponseWriter, req *http.Request) {
	//getting chirps
	var tagFixed []Cherp
	arrayOfChirps, err := conf.qs.GetAllChirps(context.Background())
	if err != nil {
		fmt.Printf("ERROR getting chirps: %s\n", err)
		os.Exit(1)
	}

	for _, entry := range arrayOfChirps {
		new := Cherp{
			ID:        entry.ID,
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
			Body:      entry.Body,
			User_id:   entry.UserID,
		}

		tagFixed = append(tagFixed, new)
	}
	//to []byte for response
	dataArray, err := json.Marshal(tagFixed)
	if err != nil {
		fmt.Printf("ERROR marshaling: %s\n", err)
		os.Exit(1)
	}
	//set response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(200)
	writer.Write(dataArray)
}

func (conf apiConfig) getSingleChirp(writer http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("chirpID")
	if idString == "" {
		errorResponse(writer, 404, "Resource not found")
	} else {
		ego, err := uuid.Parse(idString)
		if err != nil {
			errorResponse(writer, 418, "Teapot can't parse uuid string")
		}
		chirpData, err := conf.qs.GetOneChirp(context.Background(), ego)
		if err != nil {
			errorResponse(writer, 404, "chirp not found")
		}
		newer := Cherp{
			ID:        chirpData.ID,
			CreatedAt: chirpData.CreatedAt,
			UpdatedAt: chirpData.UpdatedAt,
			Body:      chirpData.Body,
			User_id:   chirpData.UserID,
		}
		jsonResponse(writer, 200, newer)
	}
}
