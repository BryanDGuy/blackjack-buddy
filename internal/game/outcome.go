package game

import (
	"fmt"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func InitialOutcomes(hands [][]card.Card) []string {
	outcomes := make([]string, len(hands))
	for i, cards := range hands {
		if hand.NewHand(cards).IsBust() {
			outcomes[i] = "Bust"
		}
	}
	return outcomes
}

func FormatOutcomeSummary(outcomes []string) string {
	if len(outcomes) == 0 {
		return ""
	}

	parts := make([]string, 0, len(outcomes))
	for i, outcome := range outcomes {
		label := outcome
		if label == "" {
			label = "Pending"
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}
	return strings.Join(parts, " | ")
}
