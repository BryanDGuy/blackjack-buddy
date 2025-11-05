package strategy

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Decision int

const (
	Hit Decision = iota
	Stand
	DoubleDown
	Split
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
	default:
		return "UNKNOWN"
	}
}

type Strategy interface {
	GetDecision(playerHand, dealerHand *hand.Hand) Decision
}

type Advisor struct {
	strategy Strategy
}

func NewAdvisor(strategy Strategy) *Advisor {
	return &Advisor{
		strategy: strategy,
	}
}

func (a *Advisor) MakeDecision(playerHand, dealerHand *hand.Hand) (Decision, error) {
	if playerHand.IsEmpty() {
		return Hit, fmt.Errorf("player hand is empty")
	}

	if dealerHand.IsEmpty() {
		return Hit, fmt.Errorf("dealer hand is empty")
	}

	if playerHand.IsBlackjack() {
		return Stand, nil
	}

	if playerHand.IsBust() {
		return Stand, nil
	}

	return a.strategy.GetDecision(playerHand, dealerHand), nil
}
