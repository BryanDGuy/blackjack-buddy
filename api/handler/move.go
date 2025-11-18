package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"maps"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/api/store"
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
		session, exists := store.Get(sessionID)
		if !exists {
			writeError(w, http.StatusNotFound, "SESSION_NOT_FOUND", "Session not found")
			return
		}

		if session.RoundState() != game.RoundStateActive {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active round")
			return
		}

		if session.Player() == nil || session.Player().ActiveHand() == nil {
			writeError(w, http.StatusConflict, "NO_ACTIVE_ROUND", "No active hand")
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

		if !isValidMove(session.Player().ActiveHand(), decision) {
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Move is not valid for current hand state")
			return
		}

		switch decision {
		case strategy.Hit:
			err = session.Player().Hit(session)
		case strategy.Stand:
			err = session.Player().Stand(session)
		case strategy.DoubleDown:
			err = session.Player().Double(session)
		case strategy.Split:
			err = session.Player().Split(session)
		default:
			writeError(w, http.StatusBadRequest, "INVALID_MOVE", "Unsupported move")
			return
		}

		if err != nil {
			if errors.Is(err, game.ErrInvalidMove) || errors.Is(err, game.ErrInvalidSplit) {
				writeError(w, http.StatusBadRequest, "INVALID_MOVE", err.Error())
				return
			}
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to apply decision")
			return
		}

		roundComplete := !session.Player().CanMove() && len(session.Player().InactiveHands()) == 0

		if roundComplete {
			if session.Dealer() != nil && session.Dealer().Hand() != nil && session.Dealer().Hand().Count() > 0 {
				for session.Dealer().Hand().Value() < 17 {
					session.Dealer().Hand().AddCard(session.DrawCard())
				}
				outcomes := session.DetermineOutcome()
				session.UpdateOutcomes(outcomes)
			}
			session.SetRoundState(game.RoundStateComplete)
		} else {
			session.SetRoundState(game.RoundStateActive)
		}

		s := session.Shoe()
		counts := make(map[string]int)
		maps.Copy(counts, s.RankCounts())

		activeHand := []string{}
		if session.Player() != nil && session.Player().ActiveHand() != nil {
			activeHand = helpers.CardsToStrings(session.Player().ActiveHand().Cards())
		}

		inactiveHands := [][]string{}
		if session.Player() != nil {
			inactiveHandsList := session.Player().InactiveHands()
			inactiveHands = make([][]string, len(inactiveHandsList))
			for i, h := range inactiveHandsList {
				inactiveHands[i] = helpers.CardsToStrings(h.Cards())
			}
		}

		completedHands := [][]string{}
		if session.Player() != nil {
			completedHandsList := session.Player().CompletedHands()
			completedHands = make([][]string, len(completedHandsList))
			for i, h := range completedHandsList {
				completedHands[i] = helpers.CardsToStrings(h.Cards())
			}
		}

		dealerCards := []string{}
		if session.Dealer() != nil && session.Dealer().Hand() != nil {
			dealerCards = helpers.CardsToStrings(session.Dealer().Hand().Cards())
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
			RoundState:     string(session.RoundState()),
			ActiveHand:     activeHand,
			InactiveHands:  inactiveHands,
			CompletedHands: completedHands,
			DealerCards:    dealerCards,
			Outcomes:       session.Outcomes(),
			ShoeState: struct {
				TotalCards int            `json:"totalCards"`
				RankCounts map[string]int `json:"rankCounts"`
			}{
				TotalCards: s.TotalCards(),
				RankCounts: counts,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

func isValidMove(playerHand *hand.Hand, decision strategy.Decision) bool {
	switch decision {
	case strategy.Hit:
		return !playerHand.IsBust() && playerHand.Value() < 21
	case strategy.Stand:
		return !playerHand.IsBust()
	case strategy.DoubleDown:
		return playerHand.Count() == 2 && !playerHand.IsBust()
	case strategy.Split:
		return playerHand.CanSplit()
	default:
		return false
	}
}
