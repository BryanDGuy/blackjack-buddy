package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
)

func NewDeal(store *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			PlayerCards []string `json:"playerCards"`
			DealerCard  string   `json:"dealerCard"`
		}
		success := false
		if !store.WithGame(r.PathValue("id"), func(g *game.Game) {
			if g.RoundState == game.RoundStateActive {
				writeError(w, http.StatusConflict, "ROUND_ALREADY_ACTIVE", "Round is already active")
				return
			}
			g.StartRound()
			response.PlayerCards = helpers.CardsToStrings(g.Player.ActiveHand.Cards)
			response.DealerCard = g.Dealer.Hand.Cards[0].ToString()
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
