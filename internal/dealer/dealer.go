package dealer

import (
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Dealer struct {
	Hand *hand.Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: nil,
	}
}
