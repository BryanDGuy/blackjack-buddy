package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func NewHint(sessions *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			Hint string `json:"hint"`
		}
		success := false
		if !sessions.WithGame(r.PathValue("id"), func(g *game.Game) {
			if g.RoundState != game.RoundStateActive {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
				return
			}
			if g.Player == nil || g.Player.ActiveHand == nil {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
				return
			}
			if g.Dealer == nil || g.Dealer.Hand == nil {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No dealer hand")
				return
			}
			decision, err := strategy.MakeDecision(g.Player.ActiveHand, g.Dealer.Hand)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to get hint")
				return
			}
			response.Hint = string(decision)
			success = true
		}) {
			writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
			return
		}
		if !success {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
