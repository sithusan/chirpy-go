package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	// to control the response. api resource concept. might worth to group the response together.
	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	const MIN_EMAIL_LENGH = 10
	const MAX_EMAIL_LENGTH = 255

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	if err := decoder.Decode(&params); err != nil {
		errorResponse(w, http.StatusInternalServerError, "Could not decode params", err)
		return
	}

	if len(params.Email) < MIN_EMAIL_LENGH {
		errorResponse(w, http.StatusBadRequest, "Email too short", nil)
		return
	}

	if len(params.Email) > MAX_EMAIL_LENGTH {
		errorResponse(w, http.StatusBadRequest, "Email too long", nil)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Cannot create user", err)
		return
	}

	jsonResponse(w, http.StatusCreated, UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) resetUserHandler(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		errorResponse(w, http.StatusForbidden, "Not allow to reset", nil)
		return
	}

	err := cfg.db.DeleteUsers(r.Context())

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "Cannot delete users", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
