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

func (h *Hand) Value() int {
	total := 0
	aces := 0

	for _, c := range h.Cards {
		if c.Rank == card.Ace {
			aces++
		} else {
			total += c.Value()
		}
	}

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

	for _, c := range h.Cards {
		if c.Rank == card.Ace {
			aces++
		} else {
			total += c.Value()
		}
	}

	if aces > 0 && total+11+aces-1 <= 21 {
		return true
	}

	return false
}

func (h *Hand) CanSplit() bool {
	if len(h.Cards) != 2 {
		return false
	}

	return h.Cards[0].Rank == h.Cards[1].Rank
}

type HandType int

const (
	Pair2 HandType = iota
	Pair3
	Pair4
	Pair5
	Pair6
	Pair7
	Pair8
	Pair9
	Pair10
	PairA
	SoftA2
	SoftA3
	SoftA4
	SoftA5
	SoftA6
	SoftA7
	SoftA8
	SoftA9
	Hard4
	Hard5
	Hard6
	Hard7
	Hard8
	Hard9
	Hard10
	Hard11
	Hard12
	Hard13
	Hard14
	Hard15
	Hard16
	Hard17
	Hard18
	Hard19
	Hard20
)

func (h *Hand) GetType() HandType {
	if h.CanSplit() {
		pairCard := h.Cards[0]

		switch pairCard.Rank {
		case card.Ace:
			return PairA
		case card.Two:
			return Pair2
		case card.Three:
			return Pair3
		case card.Four:
			return Pair4
		case card.Five:
			return Pair5
		case card.Six:
			return Pair6
		case card.Seven:
			return Pair7
		case card.Eight:
			return Pair8
		case card.Nine:
			return Pair9
		default:
			return Pair10
		}
	}

	value := h.Value()

	if h.IsSoft() {
		switch value {
		case 13:
			return SoftA2
		case 14:
			return SoftA3
		case 15:
			return SoftA4
		case 16:
			return SoftA5
		case 17:
			return SoftA6
		case 18:
			return SoftA7
		case 19:
			return SoftA8
		default:
			return SoftA9
		}
	}

	switch {
	case value <= 4:
		return Hard4
	case value == 5:
		return Hard5
	case value == 6:
		return Hard6
	case value == 7:
		return Hard7
	case value == 8:
		return Hard8
	case value == 9:
		return Hard9
	case value == 10:
		return Hard10
	case value == 11:
		return Hard11
	case value == 12:
		return Hard12
	case value == 13:
		return Hard13
	case value == 14:
		return Hard14
	case value == 15:
		return Hard15
	case value == 16:
		return Hard16
	case value == 17:
		return Hard17
	case value == 18:
		return Hard18
	case value == 19:
		return Hard19
	default:
		return Hard20
	}
}
