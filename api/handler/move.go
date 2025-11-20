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
	"github.com/bryan/blackjack-buddy/internal/hand"
	playerpkg "github.com/bryan/blackjack-buddy/internal/player"
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

		if !isValidMove(player.ActiveHand, decision) {
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Move is not valid for current hand state")
			return
		}

		var err error
		switch decision {
		case strategy.Hit:
			err = g.Hit()
		case strategy.Stand:
			err = g.Stand()
		case strategy.DoubleDown:
			err = g.Double()
		case strategy.Split:
			err = g.Split()
		default:
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Unsupported move")
			return
		}

		if err != nil {
			if errors.Is(err, playerpkg.ErrInvalidMove) || errors.Is(err, playerpkg.ErrInvalidSplit) {
				writeError(w, http.StatusBadRequest, "INVALID_MOVE", err.Error())
				return
			}
			if errors.Is(err, playerpkg.ErrNoActiveHand) {
				writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to apply decision")
			return
		}

		roundComplete := !player.CanMove() && len(player.InactiveHands) == 0

		if roundComplete {
			if g.Dealer != nil && g.Dealer.Hand != nil && len(g.Dealer.Hand.Cards) > 0 {
				for g.Dealer.Hand.Value() < 17 {
					g.Dealer.Hand.AddCard(g.DrawCard())
				}
				outcomes := g.DetermineOutcome()
				g.UpdateOutcomes(outcomes)
			}
			g.RoundState = game.RoundStateComplete
		} else {
			g.RoundState = game.RoundStateActive
		}

		counts := make(map[string]int)
		maps.Copy(counts, g.Shoe.RankCounts)

		activeHand := []string{}
		if player.ActiveHand != nil {
			activeHand = helpers.CardsToStrings(player.ActiveHand.Cards)
		}

		inactiveHands := make([][]string, len(player.InactiveHands))
		for i, h := range player.InactiveHands {
			inactiveHands[i] = helpers.CardsToStrings(h.Cards)
		}

		completedHands := make([][]string, len(player.CompletedHands))
		for i, h := range player.CompletedHands {
			completedHands[i] = helpers.CardsToStrings(h.Cards)
		}

		dealerCards := []string{}
		if g.Dealer != nil && g.Dealer.Hand != nil {
			dealerCards = helpers.CardsToStrings(g.Dealer.Hand.Cards)
		}

		resp := struct {
			RoundState     string         `json:"roundState"`
			ActiveHand     []string       `json:"activeHand"`
			InactiveHands  [][]string     `json:"inactiveHands"`
			CompletedHands [][]string     `json:"completedHands"`
			DealerCards    []string       `json:"dealerCards"`
			Outcomes       []game.Outcome `json:"outcomes"`
			ShoeState      struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			} `json:"shoeState"`
		}{
			RoundState:     string(g.RoundState),
			ActiveHand:     activeHand,
			InactiveHands:  inactiveHands,
			CompletedHands: completedHands,
			DealerCards:    dealerCards,
			Outcomes:       g.Outcomes,
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

func isValidMove(playerHand *hand.Hand, decision strategy.Decision) bool {
	switch decision {
	case strategy.Hit:
		return !playerHand.IsBust() && playerHand.Value() < 21
	case strategy.Stand:
		return !playerHand.IsBust()
	case strategy.DoubleDown:
		return len(playerHand.Cards) == 2 && !playerHand.IsBust()
	case strategy.Split:
		return playerHand.CanSplit()
	default:
		return false
	}
}
