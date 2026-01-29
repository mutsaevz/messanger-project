

package dto

import (
	"errors"
	"strings"
	"time"
)

// CreateMessageRequest represents request body for POST /chats/{id}/messages
type CreateMessageRequest struct {
	Text string `json:"text"`
}

// Validate validates and normalizes CreateMessageRequest
func (r *CreateMessageRequest) Validate() error {
	r.Text = strings.TrimSpace(r.Text)

	if r.Text == "" {
		return errors.New("text must not be empty")
	}

	if len(r.Text) < 1 || len(r.Text) > 5000 {
		return errors.New("text length must be between 1 and 5000 characters")
	}

	return nil
}

// MessageResponse represents message response payload
type MessageResponse struct {
	ID        uint      `json:"id"`
	ChatID    uint      `json:"chat_id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}