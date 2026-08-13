package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type moveRequest struct {
	Move string `json:"move"`
}

func NewMove(sessions *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response struct {
			RoundState      string         `json:"roundState"`
			ActiveHand      []string       `json:"activeHand"`
			UnresolvedHands [][]string     `json:"unresolvedHands"`
			ResolvedHands   [][]string     `json:"resolvedHands"`
			DealerCards     []string       `json:"dealerCards"`
			Outcomes        []game.Outcome `json:"outcomes"`
		}
		success := false
		if !sessions.WithGame(r.PathValue("id"), func(g *game.Game) {
			if g.RoundState != game.RoundStateActive {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
				return
			}
			if g.ActiveHand == nil {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
				return
			}
			var req moveRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
				return
			}
			if err := g.ApplyMove(strategy.Decision(req.Move)); err != nil {
				if errors.Is(err, game.ErrInvalidMove) || errors.Is(err, game.ErrInvalidSplit) {
					writeError(w, http.StatusBadRequest, "INVALID_MOVE", err.Error())
					return
				}
				if errors.Is(err, game.ErrNoActiveHand) {
					writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", err.Error())
					return
				}
				writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to apply decision")
				return
			}
			response.RoundState = string(g.RoundState)
			response.ActiveHand = []string{}
			if g.ActiveHand != nil {
				response.ActiveHand = cardsToStrings(g.ActiveHand.Cards)
			}
			response.UnresolvedHands = make([][]string, len(g.UnresolvedHands))
			for i, h := range g.UnresolvedHands {
				response.UnresolvedHands[i] = cardsToStrings(h.Cards)
			}
			response.ResolvedHands = make([][]string, len(g.ResolvedHands))
			for i, h := range g.ResolvedHands {
				response.ResolvedHands[i] = cardsToStrings(h.Cards)
			}
			response.DealerCards = []string{}
			if g.DealerHand != nil {
				dealerCards := g.DealerHand.Cards
				if g.RoundState == game.RoundStateActive && len(dealerCards) > 0 {
					dealerCards = dealerCards[:1]
				}
				response.DealerCards = cardsToStrings(dealerCards)
			}
			response.Outcomes = append([]game.Outcome{}, g.Outcomes...)
			success = true
		}) {
			writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
			return
		}
		if !success {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
