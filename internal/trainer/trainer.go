package trainer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type Trainer struct {
	advisor *strategy.Advisor
	rng     *rand.Rand
}

func NewTrainer(strat strategy.Strategy) *Trainer {
	return &Trainer{
		advisor: strategy.NewAdvisor(strat),
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

type Scenario struct {
	PlayerCards []string `json:"playerCards"`
	DealerCard  string   `json:"dealerCard"`
}

type ScenarioRequest struct {
	SkipTrivial bool `json:"skipTrivial"`
}

type CheckRequest struct {
	PlayerCards    []string   `json:"playerCards"`
	DealerCard     string     `json:"dealerCard"`
	Decision       string     `json:"decision"`
	QueuedHands    [][]string `json:"queuedHands"`
	CompletedHands [][]string `json:"completedHands"`
}

type CheckResponse struct {
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

func (t *Trainer) generateScenario(skipTrivial bool) Scenario {
	for {
		playerCards := []card.Card{t.randomCard(), t.randomCard()}
		dealerCard := t.randomCard()

		if skipTrivial && t.isTrivialHand(playerCards) {
			continue
		}

		return Scenario{
			PlayerCards: []string{playerCards[0].ToString(), playerCards[1].ToString()},
			DealerCard:  dealerCard.ToString(),
		}
	}
}

func (t *Trainer) isTrivialHand(playerCards []card.Card) bool {
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

func (t *Trainer) randomCard() card.Card {
	ranks := []card.Rank{
		card.Ace, card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King,
	}
	return card.NewCard(ranks[t.rng.Intn(len(ranks))])
}

func (t *Trainer) parseDecision(dec string) (strategy.Decision, error) {
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

func (t *Trainer) parseCard(s string) (card.Card, error) {
	var rank card.Rank
	switch s {
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
		return card.Card{}, fmt.Errorf("invalid card: %s", s)
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

func (t *Trainer) cardsFromStrings(values []string) ([]card.Card, error) {
	cards := make([]card.Card, len(values))
	for i, s := range values {
		c, err := t.parseCard(s)
		if err != nil {
			return nil, err
		}
		cards[i] = c
	}
	return cards, nil
}

func (t *Trainer) handsFromStrings(data [][]string) ([][]card.Card, error) {
	hands := make([][]card.Card, len(data))
	for i, seq := range data {
		cards, err := t.cardsFromStrings(seq)
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

func (t *Trainer) evaluateAllHands(hands [][]card.Card, dealerHand *hand.Hand) ([]card.Card, []string) {
	dealerCards := t.finishDealer(dealerHand.Cards)
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

func (t *Trainer) finishDealer(cards []card.Card) []card.Card {
	dealerHand := hand.NewHand(cards)
	if len(dealerHand.Cards) == 1 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, t.randomCard()))
	}

	for dealerHand.Value() < 17 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, t.randomCard()))
	}

	return dealerHand.Cards
}

func (t *Trainer) handleScenario(w http.ResponseWriter, r *http.Request) {
	var req ScenarioRequest
	if r.Method == "POST" {
		json.NewDecoder(r.Body).Decode(&req)
	}
	scenario := t.generateScenario(req.SkipTrivial)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)
}

func (t *Trainer) handleCheck(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playerCards, err := t.cardsFromStrings(req.PlayerCards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dealerCard, err := t.parseCard(req.DealerCard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queuedHands, err := t.handsFromStrings(req.QueuedHands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	completedHands, err := t.handsFromStrings(req.CompletedHands)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playerHand := hand.NewHand(playerCards)
	dealerHand := hand.NewHand([]card.Card{dealerCard})

	correctDecision, err := t.advisor.MakeDecision(playerHand, dealerHand)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userDecision, err := t.parseDecision(req.Decision)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := CheckResponse{
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
		playerCards = append(playerCards, t.randomCard())
		playerHand = hand.NewHand(playerCards)
		response.PlayerCards = cardsToStrings(playerCards)
		if playerHand.IsBust() || playerHand.Value() == 21 {
			response.RoundComplete = true
		}
		if playerHand.IsBust() {
			response.Outcome = "Bust"
		}
	case strategy.DoubleDown:
		playerCards = append(playerCards, t.randomCard())
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
		firstHand := []card.Card{first, t.randomCard()}
		secondHand := []card.Card{second, t.randomCard()}
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
			dealerCards, outcomes := t.evaluateAllHands(completedHands, dealerHand)
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

func (t *Trainer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlContent))
}

func (t *Trainer) Start(port int) error {
	http.HandleFunc("/", t.handleIndex)
	http.HandleFunc("/api/scenario", t.handleScenario)
	http.HandleFunc("/api/check", t.handleCheck)

	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("http://localhost%s\n", addr)
	return http.ListenAndServe(addr, nil)
}

const htmlContent = `<!DOCTYPE html>
<html>
<head>
<title>Blackjack Trainer</title>
<style>
* { margin: 0; padding: 0; box-sizing: border-box; }
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #1a1a1a; color: #fff; margin: 0; }
.layout { min-height: 100vh; display: flex; justify-content: center; align-items: center; }
#app { width: 100%; max-width: 600px; padding: 20px; position: relative; }
.card { width: 80px; height: 112px; background: #fff; border-radius: 8px; margin: 0 8px; font-size: 32px; color: #000; display: flex; align-items: center; justify-content: center; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.3); }
.cards { text-align: center; margin: 16px 0; display: flex; justify-content: center; align-items: center; flex-wrap: nowrap; }
.label { font-size: 14px; color: #888; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 1px; text-align: center; }
.section { margin: 32px 0; }
.buttons { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
button { padding: 16px; background: #333; border: 2px solid #555; border-radius: 8px; color: #fff; font-size: 16px; cursor: pointer; transition: all 0.2s; }
button:hover { background: #444; border-color: #777; }
button:active { transform: scale(0.98); }
.result { padding: 20px; border-radius: 8px; text-align: center; font-size: 18px; margin-top: 20px; min-height: 64px; display: flex; align-items: center; justify-content: center; }
.correct { background: #1a472a; border: 2px solid #2d8659; }
.incorrect { background: #4a1a1a; border: 2px solid #862d2d; }
.outcome-box { background: #222; border: 2px solid #555; min-height: 64px; display: flex; align-items: center; justify-content: center; margin-top: 12px; }
.next-btn { margin-top: 20px; width: 100%; padding: 14px; background: #2d2d2d; color: #fff; border: 2px solid #555; border-radius: 8px; font-size: 16px; cursor: pointer; visibility: hidden; pointer-events: none; }
.next-btn.visible { visibility: visible; pointer-events: auto; }
.next-btn:hover { background: #3a3a3a; }
.stats { display: flex; justify-content: space-between; padding: 16px; background: #222; border-radius: 8px; margin-bottom: 24px; }
.stat { text-align: center; }
.stat-value { font-size: 28px; font-weight: bold; }
.stat-label { font-size: 12px; color: #888; margin-top: 4px; }
.options { text-align: center; margin: 16px 0; }
.checkbox { font-size: 14px; color: #fff; cursor: pointer; }
.checkbox input { margin-right: 8px; cursor: pointer; }
.info-wrap { position: fixed; top: 16px; right: 16px; display: flex; gap: 12px; align-items: flex-start; flex-direction: row-reverse; z-index: 10; }
.info-button { width: 36px; height: 36px; border-radius: 50%; border: 2px solid #555; background: rgba(0,0,0,0.4); color: #fff; display: flex; align-items: center; justify-content: center; cursor: pointer; font-weight: bold; box-shadow: 0 4px 12px rgba(0,0,0,0.4); }
.info-panel { width: 320px; background: rgba(0,0,0,0.92); border: 2px solid #555; border-radius: 12px; padding: 16px; font-size: 12px; line-height: 1.4; opacity: 0; pointer-events: none; transition: opacity 0.2s ease; }
.info-panel table { width: 100%; border-collapse: collapse; margin-top: 12px; }
.info-panel th, .info-panel td { padding: 4px; text-align: center; border: 1px solid rgba(255,255,255,0.1); font-size: 11px; }
.info-panel .section { font-weight: bold; text-align: left; padding: 6px 0 4px; }
.info-wrap:hover .info-panel { opacity: 1; pointer-events: auto; }
.cell-stand { background: #7b1f1f; }
.cell-hit { background: #2d4f2d; }
.cell-double { background: #7b6b1f; }
.cell-split { background: #1f3f7b; }
.cell-stand, .cell-hit, .cell-double, .cell-split { color: #fff; }
</style>
</head>
<body>
<div class="info-wrap">
<div class="info-button">i</div>
<div class="info-panel">
<div class="section">Hard Hands</div>
<table id="hard-table"></table>
<div class="section">Soft Hands</div>
<table id="soft-table"></table>
<div class="section">Pairs</div>
<table id="pair-table"></table>
</div>
</div>
<div class="layout">
<div id="app">
<div class="stats">
<div class="stat"><div class="stat-value" id="correct">0</div><div class="stat-label">CORRECT</div></div>
<div class="stat"><div class="stat-value" id="total">0</div><div class="stat-label">TOTAL</div></div>
<div class="stat"><div class="stat-value" id="percent">0%</div><div class="stat-label">ACCURACY</div></div>
</div>
<div class="options">
<label class="checkbox"><input type="checkbox" id="skip-trivial" checked> Skip Trivial</label>
</div>
<div class="section">
<div class="label">Dealer Card</div>
<div class="cards" id="dealer-card"></div>
</div>
<div class="section">
<div class="label">Your Cards</div>
<div class="cards" id="player-cards"></div>
</div>
<div class="buttons">
<button onclick="decide('HIT')">HIT (H)</button>
<button onclick="decide('STAND')">STAND (S)</button>
<button onclick="decide('DOUBLE DOWN')">DOUBLE DOWN (D)</button>
<button onclick="decide('SPLIT')">SPLIT (P)</button>
</div>
<div id="result"></div>
<div id="outcome"></div>
<button id="next" class="next-btn" onclick="loadScenario()">Next (N)</button>
</div>
</div>
<script>
let currentScenario = null;
let stats = { correct: 0, total: 0 };
const KEY_BINDINGS = {
    h: 'HIT',
    s: 'STAND',
    d: 'DOUBLE DOWN',
    p: 'SPLIT',
    n: 'NEXT'
};

const TABLE_HEADERS = ['', 'A', '10', '9', '8', '7', '6', '5', '4', '3', '2'];

const HARD_MATRIX = [
    ['17+', 'S','S','S','S','S','S','S','S','S','S'],
    ['16',  'H','H','H','H','H','S','S','S','S','S'],
    ['15',  'H','H','H','H','H','S','S','S','S','S'],
    ['14',  'H','H','H','H','H','S','S','S','S','S'],
    ['13',  'H','H','H','H','H','S','S','S','S','S'],
    ['12',  'H','H','H','H','H','S','S','S','H','H'],
    ['11',  'D','D','D','D','D','D','D','D','D','D'],
    ['10',  'H','H','D','D','D','D','D','D','D','D'],
    ['9',   'H','H','H','H','H','D','D','D','D','H'],
    ['8-',  'H','H','H','H','H','H','H','H','H','H'],
];

const SOFT_MATRIX = [
    ['A9', 'S','S','S','S','S','S','S','S','S','S'],
    ['A8', 'S','S','S','S','S','S','S','S','S','S'],
    ['A7', 'H','H','H','S','S','D','D','D','D','S'],
    ['A6', 'H','H','H','H','H','D','D','D','D','H'],
    ['A5', 'H','H','H','H','H','D','D','D','H','H'],
    ['A4', 'H','H','H','H','H','D','D','D','H','H'],
    ['A3', 'H','H','H','H','H','D','D','H','H','H'],
    ['A2', 'H','H','H','H','H','D','D','H','H','H'],
];

const PAIR_MATRIX = [
    ['AA',    'SP','SP','SP','SP','SP','SP','SP','SP','SP','SP'],
    ['10 10', 'S','S','S','S','S','S','S','S','S','S'],
    ['9 9',   'S','S','SP','SP','S','SP','SP','SP','SP','SP'],
    ['8 8',   'SP','SP','SP','SP','SP','SP','SP','SP','SP','SP'],
    ['7 7',   'H','H','H','H','SP','SP','SP','SP','SP','SP'],
    ['6 6',   'H','H','H','H','H','SP','SP','SP','SP','SP'],
    ['5 5',   'H','H','D','D','D','D','D','D','D','D'],
    ['4 4',   'H','H','H','H','H','SP','SP','H','H','H'],
    ['3 3',   'H','H','H','H','SP','SP','SP','SP','SP','SP'],
    ['2 2',   'H','H','H','H','SP','SP','SP','SP','SP','SP'],
];

function renderScenario() {
    if (!currentScenario) return;
    document.getElementById('player-cards').innerHTML = currentScenario.playerCards.map(c => '<div class="card">' + c + '</div>').join('');
    const dealerCards = Array.isArray(currentScenario.dealerCards) && currentScenario.dealerCards.length
        ? currentScenario.dealerCards
        : [currentScenario.dealerCard];
    document.getElementById('dealer-card').innerHTML = dealerCards.map(c => '<div class="card">' + c + '</div>').join('');
}

function renderTable(targetId, matrix) {
    const target = document.getElementById(targetId);
    if (!target) return;
    const headerRow = '<tr>' + TABLE_HEADERS.map(h => '<th>' + h + '</th>').join('') + '</tr>';
    const body = matrix.map(row => {
        const [label, ...cells] = row;
        const cellHtml = cells.map(val => '<td class="' + classForDecision(val) + '">' + val + '</td>').join('');
        return '<tr><td>' + label + '</td>' + cellHtml + '</tr>';
    }).join('');
    target.innerHTML = headerRow + body;
}

function classForDecision(val) {
    switch (val) {
        case 'S': return 'cell-stand';
        case 'H': return 'cell-hit';
        case 'D': return 'cell-double';
        case 'SP': return 'cell-split';
        default: return '';
    }
}

async function loadScenario() {
    const skipTrivial = document.getElementById('skip-trivial').checked;
    const res = await fetch('/api/scenario', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ skipTrivial: skipTrivial })
    });
    currentScenario = await res.json();
    currentScenario.dealerCards = [currentScenario.dealerCard];
    currentScenario.queuedHands = [];
    currentScenario.completedHands = [];
    currentScenario.completedOutcomes = [];
    renderScenario();
    const resultDiv = document.getElementById('result');
    resultDiv.className = 'result';
    resultDiv.innerHTML = 'Result: -';
    const outcomeDiv = document.getElementById('outcome');
    outcomeDiv.className = 'result outcome-box';
    outcomeDiv.innerHTML = 'Outcome: -';
    const nextBtn = document.getElementById('next');
    nextBtn.classList.remove('visible');
}

async function decide(decision) {
    if (!currentScenario) return;
    const nextBtn = document.getElementById('next');
    const res = await fetch('/api/check', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            playerCards: currentScenario.playerCards,
            dealerCard: currentScenario.dealerCard,
            decision: decision,
            queuedHands: currentScenario.queuedHands || [],
            completedHands: currentScenario.completedHands || []
        })
    });
    const result = await res.json();

    stats.total++;
    if (result.correct) stats.correct++;

    updateStats();

    const resultDiv = document.getElementById('result');
    const outcomeDiv = document.getElementById('outcome');
    if (result.correct) {
        resultDiv.className = 'result correct';
        resultDiv.innerHTML = 'Correct: ' + result.correctDecision;
    } else {
        resultDiv.className = 'result incorrect';
        resultDiv.innerHTML = 'Incorrect: ' + result.correctDecision;
    }
    if (result.outcome) {
        outcomeDiv.className = 'result outcome-box';
        outcomeDiv.innerHTML = 'Outcome: ' + result.outcome;
    } else {
        outcomeDiv.className = 'result outcome-box';
        outcomeDiv.innerHTML = 'Outcome: -';
    }

    if (!result.correct || result.restart) {
        nextBtn.classList.add('visible');
        return;
    }

    if (Array.isArray(result.playerCards) && result.playerCards.length) {
        currentScenario.playerCards = result.playerCards;
    }
    if (Array.isArray(result.dealerCards) && result.dealerCards.length) {
        currentScenario.dealerCards = result.dealerCards;
    }
    currentScenario.queuedHands = Array.isArray(result.queuedHands) ? result.queuedHands : [];
    currentScenario.completedHands = Array.isArray(result.completedHands) ? result.completedHands : [];
    currentScenario.completedOutcomes = Array.isArray(result.completedOutcomes) ? result.completedOutcomes : [];
    renderScenario();

    if (result.roundComplete) {
        nextBtn.classList.add('visible');
    }
}

function updateStats() {
    document.getElementById('correct').textContent = stats.correct;
    document.getElementById('total').textContent = stats.total;
    const pct = stats.total > 0 ? Math.round((stats.correct / stats.total) * 100) : 0;
    document.getElementById('percent').textContent = pct + '%';
}

function handleKey(event) {
    if (event.target && (event.target.tagName === 'INPUT' || event.target.tagName === 'SELECT' || event.target.tagName === 'TEXTAREA')) {
        return;
    }
    const action = KEY_BINDINGS[event.key.toLowerCase()];
    if (!action) return;
    event.preventDefault();
    if (action === 'NEXT') {
        const nextBtn = document.getElementById('next');
        if (nextBtn.classList.contains('visible')) {
            loadScenario();
        }
        return;
    }
    decide(action);
}

renderTable('hard-table', HARD_MATRIX);
renderTable('soft-table', SOFT_MATRIX);
renderTable('pair-table', PAIR_MATRIX);
loadScenario();
document.addEventListener('keydown', handleKey);
</script>
</body>
</html>`
