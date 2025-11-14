package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bryan/blackjack-buddy/api/helpers"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type checkRequest struct {
	PlayerCards    []string   `json:"playerCards"`
	DealerCard     string     `json:"dealerCard"`
	Decision       string     `json:"decision"`
	QueuedHands    [][]string `json:"queuedHands"`
	CompletedHands [][]string `json:"completedHands"`
	Pot            int        `json:"pot"`
	Bet            int        `json:"bet"`
	TotalWinnings  int        `json:"totalWinnings"`
}

type checkResponse struct {
	Correct           bool           `json:"correct"`
	CorrectDecision   string         `json:"correctDecision"`
	UserDecision      string         `json:"userDecision"`
	PlayerCards       []string       `json:"playerCards"`
	DealerCards       []string       `json:"dealerCards"`
	QueuedHands       [][]string     `json:"queuedHands"`
	CompletedHands    [][]string     `json:"completedHands"`
	CompletedOutcomes []string       `json:"completedOutcomes"`
	Outcome           string         `json:"outcome"`
	RoundComplete     bool           `json:"roundComplete"`
	Restart           bool           `json:"restart"`
	Pot               int            `json:"pot"`
	Bet               int            `json:"bet"`
	RoundWinnings     int            `json:"roundWinnings"`
	TotalWinnings     int            `json:"totalWinnings"`
	DeckState         game.DeckState `json:"deckState"`
	Hint              string         `json:"hint"`
}

func NewCheck(advisor *strategy.Advisor, engine *game.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req checkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		playerCards, err := helpers.CardsFromStrings(req.PlayerCards)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		dealerCard, err := helpers.ParseCardFromString(req.DealerCard)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		queuedHands, err := helpers.HandsFromStrings(req.QueuedHands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		completedHands, err := helpers.HandsFromStrings(req.CompletedHands)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		playerHand := hand.NewHand(playerCards)
		dealerHand := hand.NewHand([]card.Card{dealerCard})

		correctDecision, err := advisor.MakeDecision(playerHand, dealerHand)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		userDecision, err := parseDecision(req.Decision)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		state := game.NewRoundState(playerCards, dealerCard, queuedHands, completedHands, game.InitialOutcomes(completedHands))

		bet := req.Bet
		if bet <= 0 {
			bet = game.DefaultBet
		}
		pot := req.Pot
		if pot <= 0 {
			pot = game.StartingPot
		}
		totalWinnings := req.TotalWinnings

		response := checkResponse{
			Correct:           userDecision == correctDecision,
			CorrectDecision:   correctDecision.ToString(),
			UserDecision:      userDecision.ToString(),
			PlayerCards:       helpers.CardsToStrings(playerCards),
			DealerCards:       []string{},
			QueuedHands:       helpers.HandsToStrings(queuedHands),
			CompletedHands:    helpers.HandsToStrings(completedHands),
			CompletedOutcomes: append([]string{}, state.Outcomes...),
			Outcome:           "",
			RoundComplete:     false,
			Restart:           false,
			Pot:               pot,
			Bet:               bet,
			RoundWinnings:     0,
			TotalWinnings:     totalWinnings,
			DeckState:         game.DeckState{},
			Hint:              "",
		}

		if !response.Correct {
			response.RoundComplete = true
			response.Restart = true
			response.DealerCards = helpers.CardsToStrings(dealerHand.Cards)
			handBets := make([]int, len(state.HandBets))
			copy(handBets, state.HandBets)
			for len(handBets) < len(state.Outcomes) {
				handBets = append(handBets, bet)
			}
			roundWinnings := game.CalculateWinnings(state.Outcomes, handBets)
			response.RoundWinnings = roundWinnings
			response.TotalWinnings = totalWinnings + roundWinnings
			response.Pot = pot + roundWinnings
			response.DeckState = engine.GetDeckState()
			hintDecision, _ := advisor.MakeDecision(playerHand, dealerHand)
			response.Hint = hintDecision.ToString()
			writeJSON(w, response)
			return
		}

		resolution, err := game.ApplyDecision(state, userDecision, engine, bet)
		if err != nil {
			if errors.Is(err, game.ErrInvalidSplit) {
				http.Error(w, "Invalid split", http.StatusBadRequest)
				return
			}
			http.Error(w, "Round resolution failed", http.StatusInternalServerError)
			return
		}

		response.PlayerCards = helpers.CardsToStrings(resolution.State.Player)
		response.QueuedHands = helpers.HandsToStrings(resolution.State.Queue)
		response.CompletedHands = helpers.HandsToStrings(resolution.State.Completed)
		response.CompletedOutcomes = append([]string{}, resolution.State.Outcomes...)
		response.DealerCards = helpers.CardsToStrings(resolution.DealerCards)
		response.Outcome = resolution.Outcome
		response.RoundComplete = resolution.RoundComplete

		if resolution.RoundComplete {
			handBets := make([]int, len(resolution.State.HandBets))
			copy(handBets, resolution.State.HandBets)
			for len(handBets) < len(resolution.State.Outcomes) {
				handBets = append(handBets, bet)
			}
			roundWinnings := game.CalculateWinnings(resolution.State.Outcomes, handBets)
			response.RoundWinnings = roundWinnings
			response.TotalWinnings = totalWinnings + roundWinnings
			response.Pot = pot + roundWinnings
		}

		response.DeckState = engine.GetDeckState()

		currentPlayerHand := hand.NewHand(resolution.State.Player)
		currentDealerHand := hand.NewHand([]card.Card{resolution.State.Dealer})
		hintDecision, _ := advisor.MakeDecision(currentPlayerHand, currentDealerHand)
		response.Hint = hintDecision.ToString()

		writeJSON(w, response)
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

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(value)
}
