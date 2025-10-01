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
	return &Hand{
		Cards: append([]card.Card{}, cards...),
	}
}

func (h *Hand) IsEmpty() bool {
	return len(h.Cards) == 0
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

func (h *Hand) Value() int {
	total := 0
	aces := 0

	// First pass: count all non-ace cards
	for _, c := range h.Cards {
		if c.IsAce() {
			aces++
		} else {
			total += c.Value()
		}
	}

	// Second pass: handle aces optimally
	for i := 0; i < aces; i++ {
		if total+11 <= 21 {
			total += 11
		} else {
			total += 1
		}
	}

	return total
}

func (h *Hand) IsBust() bool {
	return h.Value() > 21
}

func (h *Hand) IsBlackjack() bool {
	return len(h.Cards) == 2 && h.Value() == 21
}

func (h *Hand) IsSoft() bool {
	if h.IsEmpty() || h.IsBust() {
		return false
	}

	total := 0
	aces := 0

	// Count non-ace cards
	for _, c := range h.Cards {
		if c.IsAce() {
			aces++
		} else {
			total += c.Value()
		}
	}

	// Check if any ace is being counted as 11
	if aces > 0 && total+11+aces-1 <= 21 {
		return true
	}

	return false
}

func (h *Hand) CanDoubleDown() bool {
	return len(h.Cards) == 2
}

func (h *Hand) CanSplit() bool {
	if len(h.Cards) != 2 {
		return false
	}

	return h.Cards[0].Rank == h.Cards[1].Rank
}

func (h *Hand) FirstCard() (card.Card, bool) {
	if h.IsEmpty() {
		return card.Card{}, false
	}
	return h.Cards[0], true
}
