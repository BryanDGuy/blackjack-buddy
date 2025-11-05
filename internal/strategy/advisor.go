package strategy

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func GetDealerCardIndex(dealerCard card.Card) int {
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
	GetDecision(playerHandType hand.HandType, dealerIdx int) Decision
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

	playerHandType := playerHand.GetType()
	dealerUpCard := dealerHand.Cards[0]
	dealerIdx := GetDealerCardIndex(dealerUpCard)

	return a.strategy.GetDecision(playerHandType, dealerIdx), nil
}
