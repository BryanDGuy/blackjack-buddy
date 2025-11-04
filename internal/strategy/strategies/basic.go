package strategies

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

var matrix = [hand.Blackjack + 1][]strategy.Decision{
	hand.Pair2:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair3:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair4:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair5:     {strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit},
	hand.Pair6:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair7:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair8:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair9:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Stand, strategy.Split, strategy.Split, strategy.Stand, strategy.Stand},
	hand.Pair10:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.PairA:     {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.SoftA2:    {strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA3:    {strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA4:    {strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA5:    {strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA6:    {strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA7:    {strategy.Hit, strategy.Stand, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA8:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.DoubleDown, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.SoftA9:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard4:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard5:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard6:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard7:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard8:     {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard9:     {strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard10:    {strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit},
	hand.Hard11:    {strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown},
	hand.Hard12:    {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard13:    {strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard14:    {strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard15:    {strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Surrender},
	hand.Hard16:    {strategy.Surrender, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit, strategy.Surrender, strategy.Surrender},
	hand.Hard17:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard18:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard19:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard20:    {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Blackjack: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
}

type Basic struct{}

func NewBasic() *Basic {
	return &Basic{}
}

func (s *Basic) GetDecision(playerHand, dealerHand *hand.Hand) (strategy.Decision, error) {
	if playerHand.IsEmpty() {
		return strategy.Hit, fmt.Errorf("player hand is empty")
	}

	if dealerHand.IsEmpty() {
		return strategy.Hit, fmt.Errorf("dealer hand is empty")
	}

	if playerHand.IsBlackjack() {
		return strategy.Stand, nil
	}

	if playerHand.IsBust() {
		return strategy.Stand, nil
	}

	dealerUpCard, _ := dealerHand.FirstCard()
	dealerIdx := strategy.GetDealerCardIndex(dealerUpCard)

	handType := playerHand.GetType()
	decision := matrix[handType][dealerIdx]

	if playerHand.CanSplit() && decision != strategy.Split {
		if playerHand.IsSoft() {
			handType = hand.GetSoftHandType(playerHand)
		} else {
			handType = hand.GetHardHandType(playerHand.Value())
		}
		decision = matrix[handType][dealerIdx]
	}

	if decision == strategy.DoubleDown && !playerHand.CanDoubleDown() {
		if playerHand.IsSoft() && playerHand.Value() == 18 {
			return strategy.Stand, nil
		}
		return strategy.Hit, nil
	}

	return decision, nil
}
