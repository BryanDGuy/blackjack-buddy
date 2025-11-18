package hand

import (
	"fmt"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
)

type Hand struct {
	cards []card.Card
}

func NewHand(cards []card.Card) *Hand {
	c := make([]card.Card, len(cards))
	copy(c, cards)

	return &Hand{
		cards: c,
	}
}

func (h *Hand) Cards() []card.Card {
	c := make([]card.Card, len(h.cards))
	copy(c, h.cards)
	return c
}

func (h *Hand) AddCard(c card.Card) {
	h.cards = append(h.cards, c)
}

func (h *Hand) GetCard(index int) card.Card {
	return h.cards[index]
}

func (h *Hand) Count() int {
	return len(h.cards)
}

func (h *Hand) IsEmpty() bool {
	return len(h.cards) == 0
}

func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

func (h *Hand) IsBlackjack() bool {
	return len(h.cards) == 2 && h.Value() == 21
}

func (h *Hand) CanSplit() bool {
	if len(h.cards) != 2 {
		return false
	}

	return h.cards[0].Rank() == h.cards[1].Rank()
}

func (h *Hand) Value() int {
	total := 0
	aces := 0

	for _, c := range h.cards {
		total += c.Value()

		if c.Rank() == card.Ace {
			aces++
		}
	}

	for total > 21 && aces > 0 {
		total -= 10
		aces--
	}

	return total
}

func (h *Hand) IsSoft() bool {
	if h.IsEmpty() || h.IsBust() || h.IsBlackjack() {
		return false
	}

	totalWithoutAces := 0
	aces := 0

	for _, c := range h.cards {
		if c.Rank() == card.Ace {
			aces++
		} else {
			totalWithoutAces += c.Value()
		}
	}

	// Checks if there is at least one ace and the hand is not a bust if the ace is counted as 11.
	return aces > 0 && totalWithoutAces+11+(aces-1) <= 21
}

func (h *Hand) ToString() string {
	if h.IsEmpty() {
		return "[]"
	}

	var cards []string
	for _, c := range h.cards {
		cards = append(cards, c.ToString())
	}
	return fmt.Sprintf("[%s]", strings.Join(cards, ", "))
}
