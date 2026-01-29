

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"test-messenger/internal/dto"
	"test-messenger/internal/service"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// POST /chats/{id}/messages
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	var req dto.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.messageService.CreateMessage(
		r.Context(),
		uint(chatID),
		&req,
	)
	if err != nil {
		if errors.Is(err, service.ErrChatNotFoundForMessage) {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}