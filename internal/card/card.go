// Package card defines blackjack playing cards.
package card

type Card struct {
	rank Rank
}

func NewCard(rank Rank) Card {
	return Card{rank: rank}
}

func (c Card) Rank() Rank {
	return c.rank
}

func (c Card) Value() int {
	switch c.Rank() {
	case Two:
		return 2
	case Three:
		return 3
	case Four:
		return 4
	case Five:
		return 5
	case Six:
		return 6
	case Seven:
		return 7
	case Eight:
		return 8
	case Nine:
		return 9
	case Ten, Jack, Queen, King:
		return 10
	case Ace:
		return 11
	default:
		return 0
	}
}

func (c Card) ToString() string {
	switch c.Rank() {
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Nine:
		return "9"
	case Ten:
		return "10"
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	default:
		return "?"
	}
}
