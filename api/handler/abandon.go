package handler

import (
	"net/http"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
)

func NewAbandon(store *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		abandoned := false
		if !store.WithGame(r.PathValue("id"), func(g *game.Game) {
			if g.RoundState != game.RoundStateActive {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
				return
			}
			g.AbandonRound()
			abandoned = true
		}) {
			writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
			return
		}
		if abandoned {
			w.WriteHeader(http.StatusNoContent)
		}
	}
}
