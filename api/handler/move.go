package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
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
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid session ID in path")
			return
		}

		sessionID := pathParts[2]
		gameSession, exists := store.Get(sessionID)
		if !exists {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
			return
		}

		if gameSession.RoundState != game.RoundStateActive {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
			return
		}

		var req moveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
			return
		}

		decision, err := parseDecision(req.Move)
		if err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", err.Error())
			return
		}

		playerHand := hand.NewHand(gameSession.ActiveHand)
		dealerHand := hand.NewHand([]card.Card{gameSession.DealerCard})

		if !isValidMove(playerHand, dealerHand, decision) {
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Move is not valid for current hand state")
			return
		}

		roundState := game.NewRoundState(
			gameSession.ActiveHand,
			gameSession.DealerCard,
			gameSession.InactiveHands,
			gameSession.CompletedHands,
			gameSession.Outcomes,
		)

		resolution, err := game.ApplyDecision(roundState, decision, gameSession.Session)
		if err != nil {
			if errors.Is(err, game.ErrInvalidSplit) {
				writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Cannot split this hand")
				return
			}
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to apply decision")
			return
		}

		gameSession.InactiveHands = resolution.State.Queue
		gameSession.CompletedHands = resolution.State.Completed
		gameSession.Outcomes = resolution.State.Outcomes

		if resolution.RoundComplete {
			if len(resolution.State.Queue) == 0 {
				gameSession.ActiveHand = nil
				gameSession.DealerCards = resolution.DealerCards
				gameSession.Outcomes = resolution.State.Outcomes
				gameSession.RoundState = game.RoundStateComplete
			} else {
				gameSession.ActiveHand = resolution.State.Player
				gameSession.RoundState = game.RoundStateActive
			}
		} else {
			gameSession.ActiveHand = resolution.State.Player
			gameSession.RoundState = game.RoundStateActive
		}

		resp := struct {
			RoundState     string         `json:"roundState"`
			ActiveHand     []string       `json:"activeHand"`
			InactiveHands  [][]string     `json:"inactiveHands"`
			CompletedHands [][]string     `json:"completedHands"`
			DealerCards    []string       `json:"dealerCards"`
			Outcomes       []string       `json:"outcomes"`
			DeckState      game.DeckState `json:"deckState"`
		}{
			RoundState:     string(gameSession.RoundState),
			ActiveHand:     helpers.CardsToStrings(gameSession.ActiveHand),
			InactiveHands:  helpers.HandsToStrings(gameSession.InactiveHands),
			CompletedHands: helpers.HandsToStrings(gameSession.CompletedHands),
			DealerCards:    helpers.CardsToStrings(gameSession.DealerCards),
			Outcomes:       gameSession.Outcomes,
			DeckState:      gameSession.Session.GetDeckState(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func parseDecision(dec string) (strategy.Decision, error) {
	switch dec {
	case "HIT":
		return strategy.Hit, nil
	case "STAND":
		return strategy.Stand, nil
	case "DOUBLE DOWN":
		return strategy.DoubleDown, nil
	case "SPLIT":
		return strategy.Split, nil
	default:
		return strategy.Hit, fmt.Errorf("invalid decision: %s", dec)
	}
}

func isValidMove(playerHand *hand.Hand, dealerHand *hand.Hand, decision strategy.Decision) bool {
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

