package strategies

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type DecisionMatrix [hand.Pair2 + 1][]strategy.Decision

type Custom struct {
	decisionMatrix DecisionMatrix
}

func (c *Custom) getDealerCardIndex(dealerCard card.Card) int {
	switch dealerCard.Rank {
	case card.Ace:
		return 0
	case card.Ten, card.Jack, card.Queen, card.King:
		return 1
	case card.Nine:
		return 2
	case card.Eight:
		return 3
	case card.Seven:
		return 4
	case card.Six:
		return 5
	case card.Five:
		return 6
	case card.Four:
		return 7
	case card.Three:
		return 8
	case card.Two:
		return 9
	default:
		panic(fmt.Sprintf("invalid dealer card: %s", dealerCard.ToString()))
	}
}

func NewCustom(decisionMatrix DecisionMatrix) *Custom {
	return &Custom{
		decisionMatrix: decisionMatrix,
	}
}

func (c *Custom) GetDecision(playerHand, dealerHand *hand.Hand) strategy.Decision {
	playerHandType := playerHand.GetType()
	dealerIdx := c.getDealerCardIndex(dealerHand.Cards[0])
	return c.decisionMatrix[playerHandType][dealerIdx]
}
