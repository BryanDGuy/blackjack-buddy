package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
)

func NewDeal(store *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed")
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

		if g.RoundState == game.RoundStateActive {
			writeError(w, http.StatusConflict, "ROUND_ALREADY_ACTIVE", "Round is already active")
			return
		}

		g.StartRound()

		resp := struct {
			PlayerCards []string `json:"playerCards"`
			DealerCard  string   `json:"dealerCard"`
		}{
			PlayerCards: helpers.CardsToStrings(g.Player.ActiveHand.Cards),
			DealerCard:  g.Dealer.Hand.Cards[0].ToString(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
