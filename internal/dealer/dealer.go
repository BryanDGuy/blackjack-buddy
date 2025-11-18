package dealer

import (
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type CardDrawer interface {
	DrawCard() card.Card
}

type Dealer struct {
	Hand *hand.Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: nil,
	}
}

func (d *Dealer) Finish(drawer CardDrawer) {
	if d.Hand == nil || len(d.Hand.Cards) == 0 {
		return
	}

	if len(d.Hand.Cards) == 1 {
		d.Hand.Cards = append(d.Hand.Cards, drawer.DrawCard())
	}

	for d.Hand.Value() < 17 {
		d.Hand.Cards = append(d.Hand.Cards, drawer.DrawCard())
	}
}
