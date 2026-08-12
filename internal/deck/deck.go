// Package deck constructs card decks.
package deck

import "github.com/bryan/blackjack-buddy/internal/card"

type Deck struct {
	Cards []card.Card
}

const cardsPerDeck = 52

func NewDeck() Deck {
	cards := make([]card.Card, 0, cardsPerDeck)
	ranks := []card.Rank{
		card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King, card.Ace,
	}

	for _, rank := range ranks {
		for range 4 {
			cards = append(cards, card.NewCard(rank))
		}
	}

	return Deck{Cards: cards}
}
