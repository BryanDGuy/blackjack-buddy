package game

import (
	"fmt"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

const (
	StartingPot = 1000
	DefaultBet  = 10
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

func CalculateWinnings(outcomes []string, handBets []int) int {
	total := 0
	for i, outcome := range outcomes {
		if outcome == "" {
			continue
		}
		bet := DefaultBet
		if i < len(handBets) && handBets[i] > 0 {
			bet = handBets[i]
		}
		switch outcome {
		case "Blackjack":
			total += bet*3/2 + bet
		case "Win":
			total += bet
		case "Lose", "Bust":
			total -= bet
		case "Push":
		}
	}
	return total
}
