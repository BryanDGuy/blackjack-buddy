package strategy

import (
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

var decisionMatrix = [hand.Pair2 + 1][]Decision{
	hand.Hard21: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Hard20: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Hard19: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Hard18: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Hard17: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Hard16: {Hit, Hit, Hit, Hit, Hit, Stand, Stand, Stand, Stand, Stand},
	hand.Hard15: {Hit, Hit, Hit, Hit, Hit, Stand, Stand, Stand, Stand, Stand},
	hand.Hard14: {Hit, Hit, Hit, Hit, Hit, Stand, Stand, Stand, Stand, Stand},
	hand.Hard13: {Hit, Hit, Hit, Hit, Hit, Stand, Stand, Stand, Stand, Stand},
	hand.Hard12: {Hit, Hit, Hit, Hit, Hit, Stand, Stand, Stand, Hit, Hit},
	hand.Hard11: {DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown},
	hand.Hard10: {Hit, Hit, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown},
	hand.Hard9:  {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, DoubleDown, DoubleDown, Hit},
	hand.Hard8:  {Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit},
	hand.Hard7:  {Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit},
	hand.Hard6:  {Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit},
	hand.Hard5:  {Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit},
	hand.Hard4:  {Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit, Hit},
	hand.SoftA9: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.SoftA8: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.SoftA7: {Hit, Hit, Hit, Stand, Stand, DoubleDown, DoubleDown, DoubleDown, DoubleDown, Stand},
	hand.SoftA6: {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, DoubleDown, DoubleDown, Hit},
	hand.SoftA5: {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, DoubleDown, Hit, Hit},
	hand.SoftA4: {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, DoubleDown, Hit, Hit},
	hand.SoftA3: {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, Hit, Hit, Hit},
	hand.SoftA2: {Hit, Hit, Hit, Hit, Hit, DoubleDown, DoubleDown, Hit, Hit, Hit},
	hand.PairA:  {Split, Split, Split, Split, Split, Split, Split, Split, Split, Split},
	hand.Pair10: {Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand, Stand},
	hand.Pair9:  {Stand, Stand, Split, Split, Stand, Split, Split, Split, Split, Split},
	hand.Pair8:  {Split, Split, Split, Split, Split, Split, Split, Split, Split, Split},
	hand.Pair7:  {Hit, Hit, Hit, Hit, Split, Split, Split, Split, Split, Split},
	hand.Pair6:  {Hit, Hit, Hit, Hit, Hit, Split, Split, Split, Split, Split},
	hand.Pair5:  {Hit, Hit, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown, DoubleDown},
	hand.Pair4:  {Hit, Hit, Hit, Hit, Hit, Split, Split, Hit, Hit, Hit},
	hand.Pair3:  {Hit, Hit, Hit, Hit, Split, Split, Split, Split, Split, Split},
	hand.Pair2:  {Hit, Hit, Hit, Hit, Split, Split, Split, Split, Split, Split},
}

var dealerCardIndexes = [...]int{9, 8, 7, 6, 5, 4, 3, 2, 1, 1, 1, 1, 0}

func dealerCardIndex(dealerCard card.Card) int {
	rank := dealerCard.Rank()
	if rank < card.Two || rank > card.Ace {
		panic("invalid dealer card")
	}
	return dealerCardIndexes[rank]
}
