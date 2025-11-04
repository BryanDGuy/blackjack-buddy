package strategy

import (
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func GetDealerCardIndex(dealerCard card.Card) int {
	switch dealerCard.Rank {
	case card.Ace:
		return 0
	case card.Two:
		return 1
	case card.Three:
		return 2
	case card.Four:
		return 3
	case card.Five:
		return 4
	case card.Six:
		return 5
	case card.Seven:
		return 6
	case card.Eight:
		return 7
	case card.Nine:
		return 8
	case card.Ten, card.Jack, card.Queen, card.King:
		return 9
	default:
		return 0
	}
}

type Decision int

const (
	Hit Decision = iota
	Stand
	DoubleDown
	Split
	Surrender
)

func (d Decision) ToString() string {
	switch d {
	case Hit:
		return "HIT"
	case Stand:
		return "STAND"
	case DoubleDown:
		return "DOUBLE DOWN"
	case Split:
		return "SPLIT"
	case Surrender:
		return "SURRENDER"
	default:
		return "UNKNOWN"
	}
}

type Strategy interface {
	GetDecision(playerHand, dealerHand *hand.Hand) (Decision, error)
}

type Advisor struct {
	strategy Strategy
}

func NewAdvisor(strategy Strategy) *Advisor {
	return &Advisor{
		strategy: strategy,
	}
}

func (a *Advisor) GetDecision(playerHand, dealerHand *hand.Hand) (Decision, error) {
	return a.strategy.GetDecision(playerHand, dealerHand)
}
