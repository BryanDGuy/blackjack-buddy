// Package dealer manages the dealer hand.
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

func (d *Dealer) RefreshHand(newHand *hand.Hand) {
	d.Hand = newHand
}
