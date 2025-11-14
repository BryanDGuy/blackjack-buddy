package game

import (
	"fmt"

	"github.com/bryan/blackjack-buddy/internal/card"
)

func ParseCard(value string) (card.Card, error) {
	var rank card.Rank

	switch value {
	case "A":
		rank = card.Ace
	case "2":
		rank = card.Two
	case "3":
		rank = card.Three
	case "4":
		rank = card.Four
	case "5":
		rank = card.Five
	case "6":
		rank = card.Six
	case "7":
		rank = card.Seven
	case "8":
		rank = card.Eight
	case "9":
		rank = card.Nine
	case "10":
		rank = card.Ten
	case "J":
		rank = card.Jack
	case "Q":
		rank = card.Queen
	case "K":
		rank = card.King
	default:
		return card.Card{}, fmt.Errorf("invalid card: %s", value)
	}

	return card.NewCard(rank), nil
}

func CardsFromStrings(values []string) ([]card.Card, error) {
	cards := make([]card.Card, len(values))
	for i, value := range values {
		cardValue, err := ParseCard(value)
		if err != nil {
			return nil, err
		}
		cards[i] = cardValue
	}
	return cards, nil
}

func HandsFromStrings(data [][]string) ([][]card.Card, error) {
	hands := make([][]card.Card, len(data))
	for i, seq := range data {
		cards, err := CardsFromStrings(seq)
		if err != nil {
			return nil, err
		}
		hands[i] = cards
	}
	return hands, nil
}

func CardsToStrings(cards []card.Card) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ToString()
	}
	return out
}

func HandsToStrings(hands [][]card.Card) [][]string {
	out := make([][]string, len(hands))
	for i, h := range hands {
		out[i] = CardsToStrings(h)
	}
	return out
}
