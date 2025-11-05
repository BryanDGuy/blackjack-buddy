package strategies

import (
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type Coward struct{}

func NewCoward() *Coward {
	return &Coward{}
}

func (s *Coward) GetDecision(playerHand, dealerHand *hand.Hand) strategy.Decision {
	return strategy.Stand
}
