package strategy

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

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

type Advisor struct{}

func NewAdvisor() *Advisor {
	return &Advisor{}
}

func (a *Advisor) GetDecision(playerHand, dealerHand *hand.Hand) (Decision, error) {
	if playerHand.IsEmpty() {
		return Hit, fmt.Errorf("player hand is empty")
	}

	if dealerHand.IsEmpty() {
		return Hit, fmt.Errorf("dealer hand is empty")
	}

	// If player has blackjack, always stand
	if playerHand.IsBlackjack() {
		return Stand, nil
	}

	// If player is bust, can't make any decision
	if playerHand.IsBust() {
		return Stand, nil
	}

	playerValue := playerHand.Value()
	dealerUpCard, _ := dealerHand.FirstCard()

	if playerHand.CanSplit() {
		decision := a.getPairDecision(playerHand, dealerUpCard)
		if decision == Split {
			return Split, nil
		}
	}

	if playerHand.IsSoft() {
		return a.getSoftHandDecision(playerHand, dealerUpCard), nil
	}

	return a.getHardHandDecision(playerValue, dealerUpCard), nil
}

func (a *Advisor) getPairDecision(playerHand *hand.Hand, dealerUpCard card.Card) Decision {
	playerCard := playerHand.Cards[0]

	switch playerCard.Rank {
	case card.Ace:
		return Split
	case card.Two:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six || dealerUpCard.Rank == card.Seven {
			return Split
		}
	case card.Three:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six || dealerUpCard.Rank == card.Seven {
			return Split
		}
	case card.Six:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six {
			return Split
		}
	case card.Seven:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six || dealerUpCard.Rank == card.Seven {
			return Split
		}
	case card.Eight:
		return Split
	case card.Nine:
		if dealerUpCard.Rank != card.Seven && dealerUpCard.Rank != card.Ten &&
			dealerUpCard.Rank != card.Jack && dealerUpCard.Rank != card.Queen &&
			dealerUpCard.Rank != card.King && dealerUpCard.Rank != card.Ace {
			return Split
		}
	}

	// Default to hit for pairs that shouldn't be split
	return Hit
}

func (a *Advisor) getSoftHandDecision(playerHand *hand.Hand, dealerUpCard card.Card) Decision {
	playerValue := playerHand.Value()

	switch playerValue {
	case 13, 14, 15, 16, 17:
		if dealerUpCard.Rank == card.Five || dealerUpCard.Rank == card.Six {
			if playerHand.CanDoubleDown() {
				return DoubleDown
			}
			return Hit
		}
		return Hit
	case 18:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six {
			if playerHand.CanDoubleDown() {
				return DoubleDown
			}
			return Stand
		}
		if dealerUpCard.Rank == card.Seven || dealerUpCard.Rank == card.Eight {
			return Stand
		}
		return Hit
	case 19:
		if dealerUpCard.Rank == card.Six {
			if playerHand.CanDoubleDown() {
				return DoubleDown
			}
			return Stand
		}
		return Stand
	default:
		return Stand
	}
}

func (a *Advisor) getHardHandDecision(playerValue int, dealerUpCard card.Card) Decision {
	switch playerValue {
	case 5, 6, 7, 8:
		return Hit
	case 9:
		if dealerUpCard.Rank == card.Three || dealerUpCard.Rank == card.Four ||
			dealerUpCard.Rank == card.Five || dealerUpCard.Rank == card.Six {
			return DoubleDown
		}
		return Hit
	case 10:
		if dealerUpCard.Rank != card.Ten && dealerUpCard.Rank != card.Jack &&
			dealerUpCard.Rank != card.Queen && dealerUpCard.Rank != card.King &&
			dealerUpCard.Rank != card.Ace {
			return DoubleDown
		}
		return Hit
	case 11:
		return DoubleDown
	case 12:
		if dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six {
			return Stand
		}
		return Hit
	case 13, 14, 15, 16:
		if dealerUpCard.Rank == card.Two || dealerUpCard.Rank == card.Three ||
			dealerUpCard.Rank == card.Four || dealerUpCard.Rank == card.Five ||
			dealerUpCard.Rank == card.Six {
			return Stand
		}
		if playerValue == 16 &&
			(dealerUpCard.Rank == card.Nine || dealerUpCard.Rank == card.Ten ||
				dealerUpCard.Rank == card.Jack || dealerUpCard.Rank == card.Queen ||
				dealerUpCard.Rank == card.King || dealerUpCard.Rank == card.Ace) {
			return Surrender
		}
		if playerValue == 15 && dealerUpCard.Rank == card.Ten {
			return Surrender
		}
		return Hit
	case 17, 18, 19, 20, 21:
		return Stand
	default:
		return Stand
	}
}
