package main

import (
	"encoding/json"
	"net/http"
)

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type successResponse struct {
		Valid bool `json:"valid"`
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
		Valid: true,
	})
}
