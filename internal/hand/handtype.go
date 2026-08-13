package hand

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
)

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
		rank := h.Cards[0].Rank()
		if rank >= card.Ten && rank <= card.King {
			return Pair10
		}
		if rank == card.Ace {
			return PairA
		}
		if rank >= card.Two && rank <= card.Nine {
			return Pair9 + HandType(card.Nine-rank)
		}
	}

	value := h.Value()

	if h.IsSoft() {
		if value >= 13 && value <= 20 {
			return SoftA9 + HandType(20-value)
		}
	}

	if value >= 4 && value <= 21 {
		return Hard21 + HandType(21-value)
	}

	panic(fmt.Sprintf("unreachable: hand value %d cannot be classified", value))
}
