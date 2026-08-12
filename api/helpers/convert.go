// Package helpers contains API conversion helpers.
package helpers

import "github.com/bryan/blackjack-buddy/internal/card"

func CardsToStrings(cards []card.Card) []string {
	out := make([]string, len(cards))
	for i, c := range cards {
		out[i] = c.ToString()
	}
	return out
}
