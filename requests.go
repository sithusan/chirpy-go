package main

import "github.com/google/uuid"

type CreateUserRequest struct {
	Email string `json:"email"`
}

type CreateChirpRequest struct {
	Body   string    `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}
