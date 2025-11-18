package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func NewHint(store *store.SessionStore, advisor *strategy.Advisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only GET method is allowed")
			return
		}

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 4 {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid session ID in path")
			return
		}

		sessionID := pathParts[2]
		session, exists := store.Get(sessionID)
		if !exists {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
			return
		}

		if session.RoundState != game.RoundStateActive {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
			return
		}

		if session.Player == nil || session.Player.ActiveHand == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
			return
		}

		if session.Dealer == nil || session.Dealer.Hand == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No dealer hand")
			return
		}

		playerHand := session.Player.ActiveHand
		dealerHand := session.Dealer.Hand

		decision, err := advisor.MakeDecision(playerHand, dealerHand)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get hint")
			return
		}

		resp := struct {
			Hint string `json:"hint"`
		}{
			Hint: decision.ToString(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
