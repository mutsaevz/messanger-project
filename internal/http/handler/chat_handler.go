
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"test-messenger/internal/dto"
	"test-messenger/internal/service"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

// POST /chats
func (h *ChatHandler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateChatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	resp, err := h.chatService.CreateChat(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

// GET /chats/{id}?limit=N
func (h *ChatHandler) GetChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil {
			if v > 0 && v <= 100 {
				limit = v
			}
		}
	}

	chat, err := h.chatService.GetChat(r.Context(), uint(id), limit)
	if err != nil {
		if errors.Is(err, service.ErrChatNotFound) {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(chat)
}

// DELETE /chats/{id}
func (h *ChatHandler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid chat id", http.StatusBadRequest)
		return
	}

	err = h.chatService.DeleteChat(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, service.ErrChatNotFound) {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
