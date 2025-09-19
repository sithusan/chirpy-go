package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/sithusan/chirpy-go/internal/database"
)

func replaceProfanes(body string) string {
	profanes := map[string]struct{}{ // empty struct allocate zero memory
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	splittedOriginal := strings.Split(body, " ")
	splittedLowered := strings.Split(strings.ToLower(body), " ")

	for i, word := range splittedLowered {
		if _, ok := profanes[word]; ok {
			splittedOriginal[i] = "****"
		}
	}

	return strings.Join(splittedOriginal, " ")
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	const MIN_BODY_LENGH = 5
	const MAX_BODY_LENGTH = 140

	decoder := json.NewDecoder(r.Body)
	params := CreateChirpRequest{}

	if err := decoder.Decode(&params); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not decode params", err)
		return
	}

	if len(params.Body) < MIN_BODY_LENGH {
		errorResponse(w, http.StatusBadRequest, "Chirp too short", nil)
	}

	if len(params.Body) > MAX_BODY_LENGTH {
		errorResponse(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   replaceProfanes(params.Body),
		UserID: params.UserId,
	})

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Cannot create chirp", err)
		return
	}

	jsonResponse(w, http.StatusCreated, ChirpResource{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})

}
