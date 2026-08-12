// Package shoe manages a combined deck shoe.
package shoe

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
)

type Shoe struct {
	Cards []card.Card
}

func NewShoe(decks []deck.Deck) Shoe {
	cards := make([]card.Card, 0)
	for _, d := range decks {
		cards = append(cards, d.Cards...)
	}

	return Shoe{
		Cards: cards,
	}
}

func (s *Shoe) Shuffle(rng *rand.Rand) {
	rng.Shuffle(len(s.Cards), func(i, j int) {
		s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i]
	})
}

func (s *Shoe) Draw() card.Card {
	drawn := s.Cards[0]
	s.Cards = s.Cards[1:]
	return drawn
}
