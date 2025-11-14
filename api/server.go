package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
	frontend "github.com/bryan/blackjack-buddy/web"
)

type server struct {
	advisor *strategy.Advisor
	rng     *rand.Rand
	ui      fs.FS
}

func newServer(strat strategy.Strategy) *server {
	dist := frontend.Dist()

	return &server{
		advisor: strategy.NewAdvisor(strat),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
		ui:      dist,
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

func (s *server) generateScenario(skipTrivial bool) scenario {
	for {
		playerCards := []card.Card{s.randomCard(), s.randomCard()}
		dealerCard := s.randomCard()

		if skipTrivial && s.isTrivialHand(playerCards) {
			continue
		}

		return scenario{
			PlayerCards: []string{playerCards[0].ToString(), playerCards[1].ToString()},
			DealerCard:  dealerCard.ToString(),
		}
	}
}

func (s *server) isTrivialHand(playerCards []card.Card) bool {
	h := hand.NewHand(playerCards)
	if h.IsBlackjack() {
		return true
	}
	handType := h.GetType()
	switch handType {
	case hand.Hard20, hand.Hard19, hand.Hard18, hand.Hard17:
		return true
	case hand.Hard8, hand.Hard7, hand.Hard6, hand.Hard5, hand.Hard4:
		return true
	default:
		return false
	}
}

func (s *server) randomCard() card.Card {
	ranks := []card.Rank{
		card.Ace, card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King,
	}
	return card.NewCard(ranks[s.rng.Intn(len(ranks))])
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

func (s *server) parseCard(sv string) (card.Card, error) {
	var rank card.Rank
	switch sv {
	case "A":
		rank = card.Ace
	case "2":
		rank = card.Two
	case "3":
		rank = card.Three
	case "4":
		rank = card.Four
	case "5":
		rank = card.Five
	case "6":
		rank = card.Six
	case "7":
		rank = card.Seven
	case "8":
		rank = card.Eight
	case "9":
		rank = card.Nine
	case "10":
		rank = card.Ten
	case "J":
		rank = card.Jack
	case "Q":
		rank = card.Queen
	case "K":
		rank = card.King
	default:
		return card.Card{}, fmt.Errorf("invalid card: %s", sv)
	}
	return card.NewCard(rank), nil
}

func cardsToStrings(cards []card.Card) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ToString()
	}
	return out
}

func (s *server) cardsFromStrings(values []string) ([]card.Card, error) {
	cards := make([]card.Card, len(values))
	for i, sv := range values {
		c, err := s.parseCard(sv)
		if err != nil {
			return nil, err
		}
		cards[i] = c
	}
	return cards, nil
}

func (s *server) handsFromStrings(data [][]string) ([][]card.Card, error) {
	hands := make([][]card.Card, len(data))
	for i, seq := range data {
		cards, err := s.cardsFromStrings(seq)
		if err != nil {
			return nil, err
		}
		hands[i] = cards
	}
	return hands, nil
}

func handsToStrings(hands [][]card.Card) [][]string {
	out := make([][]string, len(hands))
	for i, h := range hands {
		out[i] = cardsToStrings(h)
	}
	return out
}

func initialOutcomes(hands [][]card.Card) []string {
	outcomes := make([]string, len(hands))
	for i, h := range hands {
		if hand.NewHand(h).IsBust() {
			outcomes[i] = "Bust"
		}
	}
	return outcomes
}

func (s *server) evaluateAllHands(hands [][]card.Card, dealerHand *hand.Hand) ([]card.Card, []string) {
	dealerCards := s.finishDealer(dealerHand.Cards)
	finalDealer := hand.NewHand(dealerCards)
	out := make([]string, len(hands))

	for i, cards := range hands {
		ph := hand.NewHand(cards)
		switch {
		case ph.IsBust():
			out[i] = "Bust"
		case len(ph.Cards) == 2 && ph.Value() == 21:
			if finalDealer.IsBlackjack() {
				out[i] = "Push"
			} else {
				out[i] = "Blackjack"
			}
		case finalDealer.IsBust():
			out[i] = "Win"
		case ph.Value() > finalDealer.Value():
			out[i] = "Win"
		case ph.Value() < finalDealer.Value():
			out[i] = "Lose"
		default:
			out[i] = "Push"
		}
	}

	return dealerCards, out
}

func formatOutcomeSummary(outcomes []string) string {
	if len(outcomes) == 0 {
		return ""
	}
	parts := make([]string, 0, len(outcomes))
	for i, outcome := range outcomes {
		label := outcome
		if label == "" {
			label = "Pending"
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}
	return strings.Join(parts, " | ")
}

func (s *server) finishDealer(cards []card.Card) []card.Card {
	dealerHand := hand.NewHand(cards)
	if len(dealerHand.Cards) == 1 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, s.randomCard()))
	}

	for dealerHand.Value() < 17 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, s.randomCard()))
	}

	return dealerHand.Cards
}

func (s *server) handleScenario(w http.ResponseWriter, r *http.Request) {
	var req scenarioRequest
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&req)
	}
	scenario := s.generateScenario(req.SkipTrivial)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)
}

func (s *server) handleCheck(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playerCards, err := s.cardsFromStrings(req.PlayerCards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dealerCard, err := s.parseCard(req.DealerCard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queuedHands, err := s.handsFromStrings(req.QueuedHands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	completedHands, err := s.handsFromStrings(req.CompletedHands)
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
		PlayerCards:     cardsToStrings(playerCards),
		QueuedHands:     handsToStrings(queuedHands),
		CompletedHands:  handsToStrings(completedHands),
	}

	completedOutcomes := initialOutcomes(completedHands)
	response.CompletedOutcomes = append([]string{}, completedOutcomes...)

	if !response.Correct {
		response.RoundComplete = true
		response.Restart = true
		response.DealerCards = cardsToStrings(dealerHand.Cards)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	switch userDecision {
	case strategy.Hit:
		playerCards = append(playerCards, s.randomCard())
		playerHand = hand.NewHand(playerCards)
		response.PlayerCards = cardsToStrings(playerCards)
		if playerHand.IsBust() || playerHand.Value() == 21 {
			response.RoundComplete = true
		}
		if playerHand.IsBust() {
			response.Outcome = "Bust"
		}
	case strategy.DoubleDown:
		playerCards = append(playerCards, s.randomCard())
		playerHand = hand.NewHand(playerCards)
		response.PlayerCards = cardsToStrings(playerCards)
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
		firstHand := []card.Card{first, s.randomCard()}
		secondHand := []card.Card{second, s.randomCard()}
		playerCards = firstHand
		queuedHands = append([][]card.Card{secondHand}, queuedHands...)
		response.PlayerCards = cardsToStrings(playerCards)
		response.QueuedHands = handsToStrings(queuedHands)
	}

	if response.RoundComplete {
		finishedHand := append([]card.Card{}, playerCards...)
		completedHands = append(completedHands, finishedHand)
		handOutcome := response.Outcome
		if handOutcome == "" {
			if hand.NewHand(finishedHand).IsBust() {
				handOutcome = "Bust"
			}
		}
		completedOutcomes = append(completedOutcomes, handOutcome)
		response.CompletedHands = handsToStrings(completedHands)
		response.CompletedOutcomes = append([]string{}, completedOutcomes...)

		if len(queuedHands) > 0 {
			next := append([]card.Card{}, queuedHands[0]...)
			queuedHands = queuedHands[1:]
			playerCards = next
			response.PlayerCards = cardsToStrings(playerCards)
			response.QueuedHands = handsToStrings(queuedHands)
			response.RoundComplete = false
			response.Outcome = ""
			response.DealerCards = cardsToStrings(dealerHand.Cards)
		} else {
			dealerCards, outcomes := s.evaluateAllHands(completedHands, dealerHand)
			response.DealerCards = cardsToStrings(dealerCards)
			for i := range completedOutcomes {
				if completedOutcomes[i] == "" && i < len(outcomes) {
					completedOutcomes[i] = outcomes[i]
				}
			}
			response.CompletedOutcomes = append([]string{}, completedOutcomes...)
			response.Outcome = formatOutcomeSummary(completedOutcomes)
		}
	} else {
		response.QueuedHands = handsToStrings(queuedHands)
		response.CompletedHands = handsToStrings(completedHands)
		response.CompletedOutcomes = append([]string{}, completedOutcomes...)
		if len(response.DealerCards) == 0 {
			response.DealerCards = cardsToStrings(dealerHand.Cards)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *server) handleUI(w http.ResponseWriter, r *http.Request) {
	if s.ui == nil {
		http.Error(w, "trainer UI not built. run `npm install` and `npm run build` inside web/", http.StatusServiceUnavailable)
		return
	}

	data, err := fs.ReadFile(s.ui, "index.html")
	if err != nil {
		http.Error(w, "trainer UI not built. run `npm install` and `npm run build` inside web/", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

func (s *server) Start(port int) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/scenario", s.handleScenario)
	mux.HandleFunc("/api/check", s.handleCheck)

	if s.ui != nil {
		fileServer := http.FileServer(http.FS(s.ui))
		mux.Handle("/assets/", fileServer)
		mux.HandleFunc("/", s.handleUI)
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "trainer UI not built. run `npm install` and `npm run build` inside web/", http.StatusServiceUnavailable)
		})
	}

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("http://localhost%s\n", addr)

	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
