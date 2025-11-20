package shoe

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
)

type Shoe struct {
	Cards      []card.Card
	TotalCards int
	RankCounts map[string]int
}

func NewShoe(decks []deck.Deck) Shoe {
	cards := make([]card.Card, 0)
	for _, d := range decks {
		cards = append(cards, d.Cards...)
	}

	counts := make(map[string]int)
	for _, c := range cards {
		rankStr := c.ToString()
		counts[rankStr]++
	}

	return Shoe{
		Cards:      cards,
		TotalCards: len(cards),
		RankCounts: counts,
	}
}

func (s *Shoe) Shuffle(rng *rand.Rand) {
	rng.Shuffle(len(s.Cards), func(i, j int) {
		s.Cards[i], s.Cards[j] = s.Cards[j], s.Cards[i]
	})
}

func (s *Shoe) Draw() card.Card {
	card := s.Cards[0]
	s.Cards = s.Cards[1:]

	rankStr := card.ToString()
	s.RankCounts[rankStr]--
	s.TotalCards--

	return card
}
