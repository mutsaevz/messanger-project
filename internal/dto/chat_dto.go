package dto

import (
	"strings"
	"errors"
)

// CreateChatRequest represents request body for POST /chats
type CreateChatRequest struct {
	Title string `json:"title"`
}

// Validate validates and normalizes CreateChatRequest
func (r *CreateChatRequest) Validate() error {
	r.Title = strings.TrimSpace(r.Title)

	if r.Title == "" {
		return errors.New("title must not be empty")
	}

	if len(r.Title) < 1 || len(r.Title) > 200 {
		return errors.New("title length must be between 1 and 200 characters")
	}

	return nil
}

// ChatResponse represents chat response payload
type ChatResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}
