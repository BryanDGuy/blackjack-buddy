package strategies

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

var strategyMatrix = [hand.Pair2 + 1][]strategy.Decision{
	hand.Hard20: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard19: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard18: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard17: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard16: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard15: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard14: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard13: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Hard12: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Hit, strategy.Hit},
	hand.Hard11: {strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown},
	hand.Hard10: {strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown},
	hand.Hard9:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit},
	hand.Hard8:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard7:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard6:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard5:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Hard4:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA9: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.SoftA8: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.SoftA7: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Stand, strategy.Stand, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Stand},
	hand.SoftA6: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit},
	hand.SoftA5: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit},
	hand.SoftA4: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit},
	hand.SoftA3: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.SoftA2: {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.PairA:  {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair10: {strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand, strategy.Stand},
	hand.Pair9:  {strategy.Stand, strategy.Stand, strategy.Split, strategy.Split, strategy.Stand, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair8:  {strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair7:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair6:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair5:  {strategy.Hit, strategy.Hit, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown, strategy.DoubleDown},
	hand.Pair4:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Hit, strategy.Hit, strategy.Hit},
	hand.Pair3:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
	hand.Pair2:  {strategy.Hit, strategy.Hit, strategy.Hit, strategy.Hit, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split, strategy.Split},
}

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

type Basic struct{}

func NewBasic() *Basic {
	return &Basic{}
}

func (s *Basic) GetDecision(playerHand, dealerHand *hand.Hand) strategy.Decision {
	playerHandType := playerHand.GetType()
	dealerIdx := GetDealerCardIndex(dealerHand.Cards[0])
	return strategyMatrix[playerHandType][dealerIdx]
}
