package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"maps"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type moveRequest struct {
	Move string `json:"move"`
}

func NewMove(store *store.SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Only POST method is allowed")
			return
		}

		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(pathParts) < 4 {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid game ID in path")
			return
		}

		gameId := pathParts[2]
		g, exists := store.Get(gameId)
		if !exists {
			writeError(w, http.StatusNotFound, "GAME_NOT_FOUND", "Game not found")
			return
		}

		if g.RoundState != game.RoundStateActive {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
			return
		}

		player := g.Player
		if player == nil || player.ActiveHand == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
			return
		}

		var req moveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		decision := strategy.Decision(req.Move)
		err := g.ApplyMove(decision)

		if err != nil {
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

		counts := make(map[string]int)
		maps.Copy(counts, g.Shoe.RankCounts)

		activeHand := []string{}
		if player.ActiveHand != nil {
			activeHand = helpers.CardsToStrings(player.ActiveHand.Cards)
		}

		unresolvedHands := make([][]string, len(player.UnresolvedHands))
		for i, h := range player.UnresolvedHands {
			unresolvedHands[i] = helpers.CardsToStrings(h.Cards)
		}

		resolvedHands := make([][]string, len(player.ResolvedHands))
		for i, h := range player.ResolvedHands {
			resolvedHands[i] = helpers.CardsToStrings(h.Cards)
		}

		dealerCards := []string{}
		if g.Dealer != nil && g.Dealer.Hand != nil {
			dealerCards = helpers.CardsToStrings(g.Dealer.Hand.Cards)
		}

		resp := struct {
			RoundState      string         `json:"roundState"`
			ActiveHand      []string       `json:"activeHand"`
			UnresolvedHands [][]string     `json:"unresolvedHands"`
			ResolvedHands   [][]string     `json:"resolvedHands"`
			DealerCards     []string       `json:"dealerCards"`
			Outcomes        []game.Outcome `json:"outcomes"`
			ShoeState       struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			} `json:"shoeState"`
		}{
			RoundState:      string(g.RoundState),
			ActiveHand:      activeHand,
			UnresolvedHands: unresolvedHands,
			ResolvedHands:   resolvedHands,
			DealerCards:     dealerCards,
			Outcomes:        g.Outcomes,
			ShoeState: struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			}{
				TotalCards: g.Shoe.TotalCards,
				RankCounts: counts,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}
