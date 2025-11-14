package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/internal/game"
)

type scenario struct {
	PlayerCards []string `json:"playerCards"`
	DealerCard  string   `json:"dealerCard"`
}

type scenarioRequest struct {
	SkipTrivial bool `json:"skipTrivial"`
}

func NewScenario(engine *game.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req scenarioRequest
		if r.Method == "POST" {
			json.NewDecoder(r.Body).Decode(&req)
		}

		sc := engine.GenerateScenario(req.SkipTrivial)
		resp := scenario{
			PlayerCards: helpers.CardsToStrings(sc.Player),
			DealerCard:  sc.Dealer.ToString(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
