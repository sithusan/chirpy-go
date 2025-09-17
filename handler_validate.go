package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not decode params", err)
		return
	}

	if len(params.Body) > 140 {
		errorResponse(w, http.StatusBadRequest, "Chirp too long", nil)
		return
	}

	jsonResponse(w, http.StatusOK, successResponse{
		CleanedBody: replaceProfanes(params.Body),
	})
}
