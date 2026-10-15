package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Mars-weald/bdv_chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (conf apiConfig) userHandler(writer http.ResponseWriter, req *http.Request) {
	//decoding request for email
	type emailGet struct {
		Email string `json:"email"`
	}

	texto := emailGet{}
	decoderPin := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoderPin.Decode(&texto)
	if err != nil {
		fmt.Printf("ERROR decoding: %s\n", err)
		os.Exit(1)
	}

	parameters := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     texto.Email,
	}

	//Create user
	userDat, err := conf.qs.CreateUser(req.Context(), parameters)
	if err != nil {
		fmt.Printf("ERROR creating user: %s\n", err)
		os.Exit(1)
	}

	//encode response
	//FUCKING TAG BULLSHIT FIX
	newUser := User{
		ID:        userDat.ID,
		CreatedAt: userDat.CreatedAt,
		UpdatedAt: userDat.UpdatedAt,
		Email:     userDat.Email,
	}
	resData, err := json.Marshal(newUser)
	if err != nil {
		fmt.Printf("ERROR marshaling json: %s\n", err)
		os.Exit(1)
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(201)
	writer.Write(resData)
}
