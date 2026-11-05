package chat

import (
	//"log/slog"
	"net/http"
	"github.com/gorilla/sessions"

	"northstar/features/chat/pages"
	"github.com/starfederation/datastar-go/datastar"
)

var messages []string

var clients []*datastar.ServerSentEventGenerator

type Handlers struct {
	sessionStore  sessions.Store
}

func NewHandlers(sessionStore sessions.Store) *Handlers {
	return &Handlers{
		sessionStore: sessionStore,
	}
}

func (h *Handlers) ChatPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.ChatPage().Render(r.Context(), w); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *Handlers) NewMessage(w http.ResponseWriter, r *http.Request) {
	type Store struct {
		Message string `json:"message"`
	}
	store := &Store{}

	if err := datastar.ReadSignals(r, store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages = append(messages, store.Message)

	// sse := datastar.NewSSE(w, r)
	//
	// err := sse.PatchElementTempl(pages.ChatMessages(messages))
	// if err != nil {
	// 	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	// }

	for _, client := range clients {
		err := client.PatchElementTempl(pages.ChatMessages(messages))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

func (h *Handlers) Messages(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	clients = append(clients, sse)

	err := sse.PatchElementTempl(pages.ChatMessages(messages))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

