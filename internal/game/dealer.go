package game

import (
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Dealer struct {
	Hand *hand.Hand
}

func NewDealer() *Dealer {
	return &Dealer{
		Hand: hand.NewHand([]card.Card{}),
	}
}

func (d *Dealer) Finish(session *Session) {
	if len(d.Hand.Cards) == 1 {
		d.Hand.Cards = append(d.Hand.Cards, session.DrawCard())
	}

	for d.Hand.Value() < 17 {
		d.Hand.Cards = append(d.Hand.Cards, session.DrawCard())
	}
}

func (d *Dealer) EvaluatePlayerHands(playerHands []*hand.Hand) []string {
	results := make([]string, len(playerHands))

	for i, playerHand := range playerHands {
		switch {
		case playerHand.IsBust():
			results[i] = "Bust"
		case len(playerHand.Cards) == 2 && playerHand.Value() == 21:
			if d.Hand.IsBlackjack() {
				results[i] = "Push"
			} else {
				results[i] = "Blackjack"
			}
		case d.Hand.IsBust():
			results[i] = "Win"
		case playerHand.Value() > d.Hand.Value():
			results[i] = "Win"
		case playerHand.Value() < d.Hand.Value():
			results[i] = "Lose"
		default:
			results[i] = "Push"
		}
	}

	return results
}
