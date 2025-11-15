package game

import (
	"maps"
	"math/rand"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/shoe"
)

const (
	decksInShoe        = 6
	reshuffleThreshold = 78
)

type Session struct {
	rng  *rand.Rand
	shoe shoe.Shoe
}

func NewSession() *Session {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	shoe := freshSession(rng)

	return &Session{
		rng:  rng,
		shoe: shoe,
	}
}

func freshSession(rng *rand.Rand) shoe.Shoe {
	decks := make([]deck.Deck, 0, decksInShoe)
	for range decksInShoe {
		d := deck.NewDeck(rng)
		d.Shuffle(rng)
		decks = append(decks, d)
	}

	s := shoe.NewShoe(decks)
	s.Shuffle(rng)

	return s
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
	if len(s.shoe.Cards) < reshuffleThreshold {
		s.shoe = freshSession(s.rng)
	}
	return s.shoe.Draw()
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

func (s *Session) GetDeckState() DeckState {
	counts := make(map[string]int)
	maps.Copy(counts, s.shoe.RankCounts)
	return DeckState{
		TotalCards: s.shoe.TotalCards,
		RankCounts: counts,
	}
}

type DeckState struct {
	TotalCards int            `json:"totalCards"`
	RankCounts map[string]int `json:"rankCounts"`
}
