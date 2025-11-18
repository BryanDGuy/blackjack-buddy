package hand

import (
	"fmt"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
)

type Hand struct {
	Cards []card.Card
}

func NewHand(cards []card.Card) *Hand {
	c := make([]card.Card, len(cards))
	copy(c, cards)

	return &Hand{
		Cards: c,
	}
}

func (h *Hand) ToString() string {
	if h.IsEmpty() {
		return "[]"
	}

	var cards []string
	for _, c := range h.Cards {
		cards = append(cards, c.ToString())
	}
	return fmt.Sprintf("[%s]", strings.Join(cards, ", "))
}

func (h *Hand) IsEmpty() bool {
	return len(h.Cards) == 0
}

func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

func (h *Hand) CanSplit() bool {
	if len(h.Cards) != 2 {
		return false
	}

	return h.Cards[0].Rank() == h.Cards[1].Rank()
}

func (h *Hand) Value() int {
	total := 0
	aces := 0

	for _, c := range h.Cards {
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

	for _, c := range h.Cards {
		if c.Rank() == card.Ace {
			aces++
		} else {
			totalWithoutAces += c.Value()
		}
	}

	// Checks if there is at least one ace and the hand is not a bust if the ace is counted as 11.
	return aces > 0 && totalWithoutAces+11+(aces-1) <= 21
}
