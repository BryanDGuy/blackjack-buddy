package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type scenario struct {
	PlayerCards []string       `json:"playerCards"`
	DealerCard  string         `json:"dealerCard"`
	Pot         int            `json:"pot"`
	Bet         int            `json:"bet"`
	DeckState   game.DeckState `json:"deckState"`
	Hint        string         `json:"hint"`
}

type scenarioRequest struct {
	SkipTrivial bool `json:"skipTrivial"`
}

func NewScenario(engine *game.Engine, advisor *strategy.Advisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req scenarioRequest
		if r.Method == "POST" {
			json.NewDecoder(r.Body).Decode(&req)
		}

		sc := engine.GenerateScenario(req.SkipTrivial)
		hint := ""
		if advisor != nil {
			playerHand := hand.NewHand(sc.Player)
			dealerHand := hand.NewHand([]card.Card{sc.Dealer})
			if decision, err := advisor.MakeDecision(playerHand, dealerHand); err == nil {
				hint = decision.ToString()
			}
		}
		resp := scenario{
			PlayerCards: helpers.CardsToStrings(sc.Player),
			DealerCard:  sc.Dealer.ToString(),
			Pot:         game.StartingPot,
			Bet:         game.DefaultBet,
			DeckState:   engine.GetDeckState(),
			Hint:        hint,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
