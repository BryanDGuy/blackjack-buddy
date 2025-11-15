package game

import (
	"math/rand"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/shoe"
	"github.com/google/uuid"
)

const (
	decksInShoe        = 6
	reshuffleThreshold = 78
)

type Session struct {
	ID   string
	Shoe shoe.Shoe
}

func NewSession() *Session {
	return &Session{
		ID:   uuid.New().String(),
		Shoe: generateShuffledShoe(),
	}
}

func generateShuffledShoe() shoe.Shoe {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
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
	if len(s.Shoe.Cards) < reshuffleThreshold {
		s.Shoe = generateShuffledShoe()
	}
	return s.Shoe.Draw()
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
