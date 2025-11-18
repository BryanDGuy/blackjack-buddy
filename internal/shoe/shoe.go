package shoe

import (
	"math/rand"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
)

type Shoe struct {
	cards      []card.Card
	totalCards int
	rankCounts map[string]int
}

func NewShoe(decks []deck.Deck) Shoe {
	cards := make([]card.Card, 0)
	for _, d := range decks {
		cards = append(cards, d.Cards()...)
	}

	counts := make(map[string]int)
	for _, c := range cards {
		rankStr := c.ToString()
		counts[rankStr]++
	}

	return Shoe{
		cards:      cards,
		totalCards: len(cards),
		rankCounts: counts,
	}
}

func (s *Shoe) Cards() []card.Card {
	c := make([]card.Card, len(s.cards))
	copy(c, s.cards)
	return c
}

func (s *Shoe) TotalCards() int {
	return s.totalCards
}

func (s *Shoe) RankCounts() map[string]int {
	counts := make(map[string]int)
	for k, v := range s.rankCounts {
		counts[k] = v
	}
	return counts
}

func (s *Shoe) Shuffle(rng *rand.Rand) {
	rng.Shuffle(len(s.cards), func(i, j int) {
		s.cards[i], s.cards[j] = s.cards[j], s.cards[i]
	})
}

func (s *Shoe) Draw() card.Card {
	card := s.cards[0]
	s.cards = s.cards[1:]

	rankStr := card.ToString()
	s.rankCounts[rankStr]--
	s.totalCards--

	return card
}
