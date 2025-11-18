package handler

import (
	"encoding/json"
	"maps"
	"net/http"
	"strings"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
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
		session, exists := store.Get(sessionID)
		if !exists {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
			return
		}

		if session.RoundState == game.RoundStateActive {
			writeError(w, http.StatusConflict, "ROUND_ALREADY_ACTIVE", "Round is already active")
			return
		}

		dealerCard := session.DrawCard()

		playerHand := hand.NewHand([]card.Card{session.DrawCard(), session.DrawCard()})
		session.Player = game.NewPlayer(playerHand)

		dealerHand := hand.NewHand([]card.Card{dealerCard})
		session.Dealer = game.NewDealer(dealerHand)

		session.RoundState = game.RoundStateActive

		s := session.Shoe
		counts := make(map[string]int)
		maps.Copy(counts, s.RankCounts)

		resp := struct {
			PlayerCards []string `json:"playerCards"`
			DealerCard  string   `json:"dealerCard"`
			ShoeState   struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			} `json:"shoeState"`
		}{
			PlayerCards: helpers.CardsToStrings(playerHand.Cards),
			DealerCard:  dealerCard.ToString(),
			ShoeState: struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			}{
				TotalCards: s.TotalCards,
				RankCounts: counts,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
