
package http

import (
	"net/http"

	"test-messenger/internal/http/handler"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter(
	chatHandler *handler.ChatHandler,
	messageHandler *handler.MessageHandler,
) *Router {
	mux := http.NewServeMux()

	// Chats
	mux.HandleFunc("POST /chats", chatHandler.CreateChat)
	mux.HandleFunc("GET /chats/{id}", chatHandler.GetChat)
	mux.HandleFunc("DELETE /chats/{id}", chatHandler.DeleteChat)

	// Messages
	mux.HandleFunc("POST /chats/{id}/messages", messageHandler.CreateMessage)

	return &Router{
		mux: mux,
	}
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
