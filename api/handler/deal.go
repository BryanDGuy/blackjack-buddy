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
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid session ID in path")
			return
		}

		sessionID := pathParts[2]
		gameSession, exists := store.Get(sessionID)
		if !exists {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
			return
		}

		if gameSession.RoundState == game.RoundStateActive {
			writeError(w, http.StatusConflict, "ROUND_ALREADY_ACTIVE", "Round is already active")
			return
		}

		deal := gameSession.Session.GenerateDeal()
		gameSession.ActiveHand = deal.Player
		gameSession.InactiveHands = nil
		gameSession.DealerCard = deal.Dealer
		gameSession.DealerCards = nil
		gameSession.Outcomes = nil
		gameSession.RoundState = game.RoundStateActive

		resp := struct {
			PlayerCards []string       `json:"playerCards"`
			DealerCard  string         `json:"dealerCard"`
			DeckState   game.DeckState `json:"deckState"`
		}{
			PlayerCards: helpers.CardsToStrings(deal.Player),
			DealerCard:  deal.Dealer.ToString(),
			DeckState:   gameSession.Session.GetDeckState(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

