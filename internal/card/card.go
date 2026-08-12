// Package card defines blackjack playing cards.
package card

type Card struct {
	rank Rank
}

var rankLabels = [...]string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A"}
var rankValues = [...]int{2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10, 11}

func NewCard(rank Rank) Card {
	return Card{rank: rank}
}

func (c Card) Rank() Rank {
	return c.rank
}

func (c Card) Value() int {
	if c.rank < Two || c.rank > Ace {
		return 0
	}
	return rankValues[c.rank]
}

func (c Card) ToString() string {
	if c.rank < Two || c.rank > Ace {
		return "?"
	}
	return rankLabels[c.rank]
}
