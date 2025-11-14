package game

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Engine struct {
	rng  *rand.Rand
	deck []card.Card
}

const (
	decksInShoe        = 6
	cardsPerDeck       = 52
	reshuffleThreshold = 78
)

func NewEngine(rng *rand.Rand) *Engine {
	e := &Engine{rng: rng}
	e.initializeDeck()
	return e
}

func (e *Engine) initializeDeck() {
	e.deck = make([]card.Card, 0, decksInShoe*cardsPerDeck)
	ranks := []card.Rank{
		card.Ace, card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King,
	}

	for i := 0; i < decksInShoe; i++ {
		for _, rank := range ranks {
			for j := 0; j < 4; j++ {
				e.deck = append(e.deck, card.NewCard(rank))
			}
		}
	}

	e.shuffle()
}

func (e *Engine) shuffle() {
	e.rng.Shuffle(len(e.deck), func(i, j int) {
		e.deck[i], e.deck[j] = e.deck[j], e.deck[i]
	})
}

type Deal struct {
	Player []card.Card
	Dealer card.Card
}

func (e *Engine) GenerateDeal(skipTrivial bool) Deal {
	for {
		player := []card.Card{e.DrawCard(), e.DrawCard()}
		dealer := e.DrawCard()

		if skipTrivial && IsTrivialHand(player) {
			continue
		}

		return Deal{
			Player: player,
			Dealer: dealer,
		}
	}
}

func (e *Engine) DrawCard() card.Card {
	if len(e.deck) < reshuffleThreshold {
		e.initializeDeck()
	}

	card := e.deck[0]
	e.deck = e.deck[1:]
	return card
}

func IsTrivialHand(cards []card.Card) bool {
	h := hand.NewHand(cards)
	if h.IsBlackjack() {
		return true
	}

	switch h.GetType() {
	case hand.Hard21, hand.Hard20, hand.Hard19, hand.Hard18, hand.Hard17:
		return true
	case hand.Hard8, hand.Hard7, hand.Hard6, hand.Hard5, hand.Hard4:
		return true
	default:
		return false
	}
}

func (e *Engine) FinishDealer(cards []card.Card) []card.Card {
	current := append([]card.Card{}, cards...)
	dealerHand := hand.NewHand(current)
	if len(dealerHand.Cards) == 1 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, e.DrawCard()))
	}

	for dealerHand.Value() < 17 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, e.DrawCard()))
	}

	return dealerHand.Cards
}

func (e *Engine) EvaluateAllHands(hands [][]card.Card, dealer card.Card) ([]card.Card, []string) {
	dealerCards := e.FinishDealer([]card.Card{dealer})
	finalDealer := hand.NewHand(dealerCards)
	results := make([]string, len(hands))

	for i, cards := range hands {
		playerHand := hand.NewHand(cards)
		switch {
		case playerHand.IsBust():
			results[i] = "Bust"
		case len(playerHand.Cards) == 2 && playerHand.Value() == 21:
			if finalDealer.IsBlackjack() {
				results[i] = "Push"
			} else {
				results[i] = "Blackjack"
			}
		case finalDealer.IsBust():
			results[i] = "Win"
		case playerHand.Value() > finalDealer.Value():
			results[i] = "Win"
		case playerHand.Value() < finalDealer.Value():
			results[i] = "Lose"
		default:
			results[i] = "Push"
		}
	}

	return dealerCards, results
}

type DeckState struct {
	TotalCards int            `json:"totalCards"`
	RankCounts map[string]int `json:"rankCounts"`
}

func (e *Engine) GetDeckState() DeckState {
	counts := make(map[string]int)
	for _, c := range e.deck {
		rankStr := c.ToString()
		counts[rankStr]++
	}
	return DeckState{
		TotalCards: len(e.deck),
		RankCounts: counts,
	}
}
