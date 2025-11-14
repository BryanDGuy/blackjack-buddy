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

type deal struct {
	PlayerCards []string       `json:"playerCards"`
	DealerCard  string         `json:"dealerCard"`
	Pot         int            `json:"pot"`
	Bet         int            `json:"bet"`
	DeckState   game.DeckState `json:"deckState"`
	Hint        string         `json:"hint"`
}

type dealRequest struct {
	SkipTrivial bool `json:"skipTrivial"`
}

func NewDeal(engine *game.Engine, advisor *strategy.Advisor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dealRequest
		if r.Method == "POST" {
			json.NewDecoder(r.Body).Decode(&req)
		}

		d := engine.GenerateDeal(req.SkipTrivial)
		playerHand := hand.NewHand(d.Player)
		dealerHand := hand.NewHand([]card.Card{d.Dealer})
		decision, _ := advisor.MakeDecision(playerHand, dealerHand)
		hint := decision.ToString()
		resp := deal{
			PlayerCards: helpers.CardsToStrings(d.Player),
			DealerCard:  d.Dealer.ToString(),
			Pot:         game.StartingPot,
			Bet:         game.DefaultBet,
			DeckState:   engine.GetDeckState(),
			Hint:        hint,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
