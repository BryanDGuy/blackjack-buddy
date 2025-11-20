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
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid game ID in path")
			return
		}

		gameId := pathParts[2]
		g, exists := store.Get(gameId)
		if !exists {
			writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
			return
		}

		if g.RoundState != game.RoundStateActive {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
			return
		}

		player := g.Player
		if player == nil || player.ActiveHand == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
			return
		}

		if g.Dealer == nil || g.Dealer.Hand == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No dealer hand")
			return
		}

		playerHand := player.ActiveHand
		dealerHand := g.Dealer.Hand

		decision, err := advisor.MakeDecision(playerHand, dealerHand)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get hint")
			return
		}

		resp := struct {
			Hint string `json:"hint"`
		}{
			Hint: string(decision),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
