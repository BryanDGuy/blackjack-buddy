package deck

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
)

type Deck struct {
	cards []card.Card
}

const cardsPerDeck = 52

func NewDeck(rng *rand.Rand) Deck {
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

	return Deck{cards: cards}
}

func (d *Deck) Cards() []card.Card {
	c := make([]card.Card, len(d.cards))
	copy(c, d.cards)
	return c
}

func (d *Deck) Count() int {
	return len(d.cards)
}

func (d *Deck) Shuffle(rng *rand.Rand) {
	rng.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}
