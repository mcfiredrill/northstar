package chat

import (
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func SetupRoutes(router chi.Router, sessionStore sessions.Store) error {
	handlers := NewHandlers(sessionStore)

	router.Get("/chat", handlers.ChatPage)
	router.Get("/messages", handlers.Messages)

	router.Post("/messages", handlers.NewMessage)
	return nil
}
