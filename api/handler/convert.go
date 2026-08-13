package handler

import "github.com/bryan/blackjack-buddy/internal/card"

func cardsToStrings(cards []card.Card) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ToString()
	}
	return out
}
