package game

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Session struct {
	rng  *rand.Rand
	deck []card.Card
}

const (
	decksInShoe        = 6
	cardsPerDeck       = 52
	reshuffleThreshold = 78
)

func NewSession(rng *rand.Rand) *Session {
	s := &Session{rng: rng}
	s.initializeDeck()
	return s
}

func (s *Session) initializeDeck() {
	s.deck = make([]card.Card, 0, decksInShoe*cardsPerDeck)
	ranks := []card.Rank{
		card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King, card.Ace,
	}

	for i := 0; i < decksInShoe; i++ {
		for _, rank := range ranks {
			for j := 0; j < 4; j++ {
				s.deck = append(s.deck, card.NewCard(rank))
			}
		}
	}

	s.shuffle()
}

func (s *Session) shuffle() {
	s.rng.Shuffle(len(s.deck), func(i, j int) {
		s.deck[i], s.deck[j] = s.deck[j], s.deck[i]
	})
}

type Deal struct {
	Player []card.Card
	Dealer card.Card
}

func (s *Session) GenerateDeal() Deal {
	player := []card.Card{s.DrawCard(), s.DrawCard()}
	dealer := s.DrawCard()

	return Deal{
		Player: player,
		Dealer: dealer,
	}
}

func (s *Session) DrawCard() card.Card {
	if len(s.deck) < reshuffleThreshold {
		s.initializeDeck()
	}

	card := s.deck[0]
	s.deck = s.deck[1:]
	return card
}

func (s *Session) FinishDealer(cards []card.Card) []card.Card {
	current := append([]card.Card{}, cards...)
	dealerHand := hand.NewHand(current)
	if len(dealerHand.Cards) == 1 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, s.DrawCard()))
	}

	for dealerHand.Value() < 17 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, s.DrawCard()))
	}

	return dealerHand.Cards
}

func (s *Session) EvaluateAllHands(hands [][]card.Card, dealer card.Card) ([]card.Card, []string) {
	dealerCards := s.FinishDealer([]card.Card{dealer})
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

func (s *Session) GetDeckState() DeckState {
	counts := make(map[string]int)
	for _, c := range s.deck {
		rankStr := c.ToString()
		counts[rankStr]++
	}
	return DeckState{
		TotalCards: len(s.deck),
		RankCounts: counts,
	}
}
