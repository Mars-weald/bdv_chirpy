package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

	err = errorResponse(writer, 400, paramaters.Body)
	if err != nil {
		return
	}

	paramaters.Cleaned_Body = censorChirp(paramaters.Body)

	paramaters.Valid = true

	//Add to database
	oozer, err := uuid.Parse(paramaters.User_id)
	if err != nil {
		fmt.Printf("ERROR parsing string to UUID")
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
		fmt.Printf("ERROR creating chirp database entry")
		return
	}

	//FUCK YOU. I FAIL BECAUSE THE FIELD NAMES ARE NOT THE EXACT. SAME. FUCK YOU.
	fixedtags := Cherp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		User_id:   chirp.UserID,
	}
	jsonResponse(writer, 201, fixedtags)
}
