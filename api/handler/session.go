package handler

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
)

func NewSession(store *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed")
			return
		}

		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		session := game.NewSession(rng)
		gameSession := game.NewGameSession(session)

		sessionID := store.Create(gameSession)

		resp := struct {
			SessionID string `json:"sessionId"`
		}{
			SessionID: sessionID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

