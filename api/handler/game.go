package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/player"
)

func NewGame(sessions *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		g := game.NewGame(player.NewPlayer(), dealer.NewDealer())

		gameID := sessions.Create(g)

		resp := struct {
			GameID string `json:"gameId"`
		}{
			GameID: gameID,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}
