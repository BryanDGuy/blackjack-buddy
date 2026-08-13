// Package shoe manages a combined deck shoe.
package shoe

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
)

type Shoe struct {
	Cards []card.Card
}

func NewShoe(decks int) Shoe {
	cards := make([]card.Card, 0, decks*52)
	for range decks {
		for rank := card.Two; rank <= card.Ace; rank++ {
			for range 4 {
				cards = append(cards, card.NewCard(rank))
			}
		}
	}
	return Shoe{Cards: cards}
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
