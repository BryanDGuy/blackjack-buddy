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

	return h.Cards[0].Rank == h.Cards[1].Rank
}

func (h *Hand) Value() int {
	total := 0
	aces := 0

	for _, c := range h.Cards {
		total += c.Value()

		if c.Rank == card.Ace {
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
		if c.Rank == card.Ace {
			aces++
		} else {
			totalWithoutAces += c.Value()
		}
	}

	// Checks if there is at least one ace and the hand is not a bust if the ace is counted as 11.
	return aces > 0 && totalWithoutAces+11+(aces-1) <= 21
}

type HandType int

const (
	Hard21 HandType = iota
	Hard20
	Hard19
	Hard18
	Hard17
	Hard16
	Hard15
	Hard14
	Hard13
	Hard12
	Hard11
	Hard10
	Hard9
	Hard8
	Hard7
	Hard6
	Hard5
	Hard4
	SoftA9
	SoftA8
	SoftA7
	SoftA6
	SoftA5
	SoftA4
	SoftA3
	SoftA2
	PairA
	Pair10
	Pair9
	Pair8
	Pair7
	Pair6
	Pair5
	Pair4
	Pair3
	Pair2
)

func (h *Hand) GetType() HandType {
	if h.CanSplit() {
		pairCard := h.Cards[0]
		switch pairCard.Rank {
		case card.Ace:
			return PairA
		case card.Ten, card.Jack, card.Queen, card.King:
			return Pair10
		case card.Nine:
			return Pair9
		case card.Eight:
			return Pair8
		case card.Seven:
			return Pair7
		case card.Six:
			return Pair6
		case card.Five:
			return Pair5
		case card.Four:
			return Pair4
		case card.Three:
			return Pair3
		case card.Two:
			return Pair2
		}
	}

	value := h.Value()

	if h.IsSoft() {
		switch value {
		case 20:
			return SoftA9
		case 19:
			return SoftA8
		case 18:
			return SoftA7
		case 17:
			return SoftA6
		case 16:
			return SoftA5
		case 15:
			return SoftA4
		case 14:
			return SoftA3
		case 13:
			return SoftA2
		}
	}

	switch value {
	case 21:
		return Hard21
	case 20:
		return Hard20
	case 19:
		return Hard19
	case 18:
		return Hard18
	case 17:
		return Hard17
	case 16:
		return Hard16
	case 15:
		return Hard15
	case 14:
		return Hard14
	case 13:
		return Hard13
	case 12:
		return Hard12
	case 11:
		return Hard11
	case 10:
		return Hard10
	case 9:
		return Hard9
	case 8:
		return Hard8
	case 7:
		return Hard7
	case 6:
		return Hard6
	case 5:
		return Hard5
	case 4:
		return Hard4
	}

	panic(fmt.Sprintf("unreachable: hand value %d cannot be classified", value))
}
