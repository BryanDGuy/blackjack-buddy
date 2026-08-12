// Package strategy provides blackjack decisions.
package strategy

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/hand"
)

type Decision string

const (
	Hit        Decision = "HIT"
	Stand      Decision = "STAND"
	DoubleDown Decision = "DOUBLE DOWN"
	Split      Decision = "SPLIT"
)

func MakeDecision(playerHand, dealerHand *hand.Hand) (Decision, error) {
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

	decision := decisionMatrix[playerHand.GetType()][dealerCardIndex(dealerHand.Cards[0])]

	if decision == DoubleDown && len(playerHand.Cards) > 2 {
		decision = Hit
	}

	return decision, nil
}
