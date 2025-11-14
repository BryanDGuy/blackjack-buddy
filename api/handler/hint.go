package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func NewHint(advisor *strategy.Advisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			PlayerCards []string `json:"playerCards"`
			DealerCard  string   `json:"dealerCard"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		playerCards, err := helpers.CardsFromStrings(req.PlayerCards)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dealerCard, err := helpers.ParseCard(req.DealerCard)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		playerHand := hand.NewHand(playerCards)
		dealerHand := hand.NewHand([]card.Card{dealerCard})

		decision, err := advisor.MakeDecision(playerHand, dealerHand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
