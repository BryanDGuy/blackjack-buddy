package dealer

import (
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Dealer struct {
	hand *hand.Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		hand: nil,
	}
}

func (d *Dealer) Hand() *hand.Hand {
	return d.hand
}

func (d *Dealer) SetHand(h *hand.Hand) {
	d.hand = h
}
