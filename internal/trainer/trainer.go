package trainer

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
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

type CheckRequest struct {
	PlayerCards []string `json:"playerCards"`
	DealerCard  string   `json:"dealerCard"`
	Decision    string   `json:"decision"`
}

type CheckResponse struct {
	Correct         bool   `json:"correct"`
	CorrectDecision string `json:"correctDecision"`
	UserDecision    string `json:"userDecision"`
}

func (t *Trainer) generateScenario() Scenario {
	playerCards := []card.Card{t.randomCard(), t.randomCard()}
	dealerCard := t.randomCard()

	return Scenario{
		PlayerCards: []string{playerCards[0].ToString(), playerCards[1].ToString()},
		DealerCard:  dealerCard.ToString(),
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

func (t *Trainer) handleScenario(w http.ResponseWriter, r *http.Request) {
	scenario := t.generateScenario()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(scenario)
}

func (t *Trainer) handleCheck(w http.ResponseWriter, r *http.Request) {
	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	playerCards := make([]card.Card, len(req.PlayerCards))
	for i, cardStr := range req.PlayerCards {
		c, err := t.parseCard(cardStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		playerCards[i] = c
	}

	dealerCard, err := t.parseCard(req.DealerCard)
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
body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #1a1a1a; color: #fff; display: flex; justify-content: center; align-items: center; min-height: 100vh; }
#app { width: 100%; max-width: 600px; padding: 20px; }
.card { width: 80px; height: 112px; background: #fff; border-radius: 8px; margin: 0 8px; font-size: 32px; color: #000; display: flex; align-items: center; justify-content: center; font-weight: bold; box-shadow: 0 4px 8px rgba(0,0,0,0.3); }
.cards { text-align: center; margin: 16px 0; display: flex; justify-content: center; align-items: center; flex-wrap: nowrap; }
.label { font-size: 14px; color: #888; margin-bottom: 12px; text-transform: uppercase; letter-spacing: 1px; text-align: center; }
.section { margin: 32px 0; }
.buttons { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
button { padding: 16px; background: #333; border: 2px solid #555; border-radius: 8px; color: #fff; font-size: 16px; cursor: pointer; transition: all 0.2s; }
button:hover { background: #444; border-color: #777; }
button:active { transform: scale(0.98); }
.result { padding: 20px; border-radius: 8px; text-align: center; font-size: 18px; margin-top: 20px; }
.correct { background: #1a472a; border: 2px solid #2d8659; }
.incorrect { background: #4a1a1a; border: 2px solid #862d2d; }
.stats { display: flex; justify-content: space-between; padding: 16px; background: #222; border-radius: 8px; margin-bottom: 24px; }
.stat { text-align: center; }
.stat-value { font-size: 28px; font-weight: bold; }
.stat-label { font-size: 12px; color: #888; margin-top: 4px; }
</style>
</head>
<body>
<div id="app">
<div class="stats">
<div class="stat"><div class="stat-value" id="correct">0</div><div class="stat-label">CORRECT</div></div>
<div class="stat"><div class="stat-value" id="total">0</div><div class="stat-label">TOTAL</div></div>
<div class="stat"><div class="stat-value" id="percent">0%</div><div class="stat-label">ACCURACY</div></div>
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
<button onclick="decide('HIT')">HIT</button>
<button onclick="decide('STAND')">STAND</button>
<button onclick="decide('DOUBLE DOWN')">DOUBLE DOWN</button>
<button onclick="decide('SPLIT')">SPLIT</button>
</div>
<div id="result"></div>
</div>
<script>
let currentScenario = null;
let stats = { correct: 0, total: 0 };

async function loadScenario() {
	const res = await fetch('/api/scenario');
	currentScenario = await res.json();
	document.getElementById('player-cards').innerHTML = currentScenario.playerCards.map(c => '<div class="card">' + c + '</div>').join('');
	document.getElementById('dealer-card').innerHTML = '<div class="card">' + currentScenario.dealerCard + '</div>';
	document.getElementById('result').innerHTML = '';
}

async function decide(decision) {
	if (!currentScenario) return;
	const res = await fetch('/api/check', {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
			playerCards: currentScenario.playerCards,
			dealerCard: currentScenario.dealerCard,
			decision: decision
		})
	});
	const result = await res.json();
	
	stats.total++;
	if (result.correct) stats.correct++;
	
	updateStats();
	
	const resultDiv = document.getElementById('result');
	if (result.correct) {
		resultDiv.className = 'result correct';
		resultDiv.innerHTML = 'Correct: ' + result.correctDecision;
	} else {
		resultDiv.className = 'result incorrect';
		resultDiv.innerHTML = 'Incorrect: ' + result.correctDecision;
	}
	
	setTimeout(loadScenario, 1500);
}

function updateStats() {
	document.getElementById('correct').textContent = stats.correct;
	document.getElementById('total').textContent = stats.total;
	const pct = stats.total > 0 ? Math.round((stats.correct / stats.total) * 100) : 0;
	document.getElementById('percent').textContent = pct + '%';
}

loadScenario();
</script>
</body>
</html>`
