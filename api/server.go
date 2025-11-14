package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type server struct {
	advisor *strategy.Advisor
	engine  *game.Engine
	ui      fs.FS
}

func newServer(strat strategy.Strategy) *server {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &server{
		advisor: strategy.NewAdvisor(strat),
		engine:  game.NewEngine(rng),
		ui:      loadUI(),
	}
}

type scenario struct {
	PlayerCards []string `json:"playerCards"`
	DealerCard  string   `json:"dealerCard"`
}

type scenarioRequest struct {
	SkipTrivial bool `json:"skipTrivial"`
}

type checkRequest struct {
	PlayerCards    []string   `json:"playerCards"`
	DealerCard     string     `json:"dealerCard"`
	Decision       string     `json:"decision"`
	QueuedHands    [][]string `json:"queuedHands"`
	CompletedHands [][]string `json:"completedHands"`
}

type checkResponse struct {
	Correct           bool       `json:"correct"`
	CorrectDecision   string     `json:"correctDecision"`
	UserDecision      string     `json:"userDecision"`
	PlayerCards       []string   `json:"playerCards"`
	DealerCards       []string   `json:"dealerCards"`
	QueuedHands       [][]string `json:"queuedHands"`
	CompletedHands    [][]string `json:"completedHands"`
	CompletedOutcomes []string   `json:"completedOutcomes"`
	Outcome           string     `json:"outcome"`
	RoundComplete     bool       `json:"roundComplete"`
	Restart           bool       `json:"restart"`
}

func (s *server) parseDecision(dec string) (strategy.Decision, error) {
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

func (s *server) handleScenario(w http.ResponseWriter, r *http.Request) {
	var req scenarioRequest
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&req)
	}

	sc := s.engine.GenerateScenario(req.SkipTrivial)
	resp := scenario{
		PlayerCards: game.CardsToStrings(sc.Player),
		DealerCard:  sc.Dealer.ToString(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *server) handleCheck(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playerCards, err := game.CardsFromStrings(req.PlayerCards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dealerCard, err := game.ParseCard(req.DealerCard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queuedHands, err := game.HandsFromStrings(req.QueuedHands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	completedHands, err := game.HandsFromStrings(req.CompletedHands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerHand := hand.NewHand(playerCards)
	dealerHand := hand.NewHand([]card.Card{dealerCard})

	correctDecision, err := s.advisor.MakeDecision(playerHand, dealerHand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDecision, err := s.parseDecision(req.Decision)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := checkResponse{
		Correct:         userDecision == correctDecision,
		CorrectDecision: correctDecision.ToString(),
		UserDecision:    userDecision.ToString(),
		PlayerCards:     game.CardsToStrings(playerCards),
		QueuedHands:     game.HandsToStrings(queuedHands),
		CompletedHands:  game.HandsToStrings(completedHands),
	}

	completedOutcomes := game.InitialOutcomes(completedHands)
	response.CompletedOutcomes = append([]string{}, completedOutcomes...)

	if !response.Correct {
		response.RoundComplete = true
		response.Restart = true
		response.DealerCards = game.CardsToStrings(dealerHand.Cards)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	switch userDecision {
	case strategy.Hit:
		playerCards = append(playerCards, s.engine.DrawCard())
		playerHand = hand.NewHand(playerCards)
		response.PlayerCards = game.CardsToStrings(playerCards)
		if playerHand.IsBust() || playerHand.Value() == 21 {
			response.RoundComplete = true
		}
		if playerHand.IsBust() {
			response.Outcome = "Bust"
		}
	case strategy.DoubleDown:
		playerCards = append(playerCards, s.engine.DrawCard())
		playerHand = hand.NewHand(playerCards)
		response.PlayerCards = game.CardsToStrings(playerCards)
		response.RoundComplete = true
		if playerHand.IsBust() {
			response.Outcome = "Bust"
		}
	case strategy.Stand:
		response.RoundComplete = true
	case strategy.Split:
		if !playerHand.CanSplit() {
			http.Error(w, "Invalid split", http.StatusBadRequest)
			return
		}
		first := playerCards[0]
		second := playerCards[1]
		firstHand := []card.Card{first, s.engine.DrawCard()}
		secondHand := []card.Card{second, s.engine.DrawCard()}
		playerCards = firstHand
		queuedHands = append([][]card.Card{secondHand}, queuedHands...)
		response.PlayerCards = game.CardsToStrings(playerCards)
		response.QueuedHands = game.HandsToStrings(queuedHands)
	}

	if response.RoundComplete {
		finishedHand := append([]card.Card{}, playerCards...)
		completedHands = append(completedHands, finishedHand)
		handOutcome := response.Outcome
		if handOutcome == "" && hand.NewHand(finishedHand).IsBust() {
			handOutcome = "Bust"
		}
		completedOutcomes = append(completedOutcomes, handOutcome)
		response.CompletedHands = game.HandsToStrings(completedHands)
		response.CompletedOutcomes = append([]string{}, completedOutcomes...)

		if len(queuedHands) > 0 {
			next := append([]card.Card{}, queuedHands[0]...)
			queuedHands = queuedHands[1:]
			playerCards = next
			response.PlayerCards = game.CardsToStrings(playerCards)
			response.QueuedHands = game.HandsToStrings(queuedHands)
			response.RoundComplete = false
			response.Outcome = ""
			response.DealerCards = game.CardsToStrings([]card.Card{dealerCard})
		} else {
			dealerCards, outcomes := s.engine.EvaluateAllHands(completedHands, dealerCard)
			response.DealerCards = game.CardsToStrings(dealerCards)
			for i := range completedOutcomes {
				if completedOutcomes[i] == "" && i < len(outcomes) {
					completedOutcomes[i] = outcomes[i]
				}
			}
			response.CompletedOutcomes = append([]string{}, completedOutcomes...)
			response.Outcome = game.FormatOutcomeSummary(completedOutcomes)
		}
	} else {
		response.QueuedHands = game.HandsToStrings(queuedHands)
		response.CompletedHands = game.HandsToStrings(completedHands)
		response.CompletedOutcomes = append([]string{}, completedOutcomes...)
		if len(response.DealerCards) == 0 {
			response.DealerCards = game.CardsToStrings([]card.Card{dealerCard})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *server) handleUI(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(s.ui, "index.html")
	if err != nil {
		http.Error(w, "trainer UI not built", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *server) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/scenario", s.handleScenario)
	mux.HandleFunc("/api/check", s.handleCheck)

	fileServer := http.FileServer(http.FS(s.ui))
	mux.Handle("/assets/", fileServer)
	mux.HandleFunc("/", s.handleUI)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("http://localhost%s\n", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
