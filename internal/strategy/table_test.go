package strategy

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func TestDecisionMatrix(t *testing.T) {
	tests := []struct {
		name         string
		playerCards  []card.Card
		dealerCard   card.Card
		wantDecision Decision
	}{
		{"hard 20 vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)}, card.NewCard(card.Ten), Stand},
		{"hard 16 vs 6", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)}, card.NewCard(card.Six), Stand},
		{"hard 16 vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)}, card.NewCard(card.Ten), Hit},
		{"hard 11 vs 5", []card.Card{card.NewCard(card.Seven), card.NewCard(card.Four)}, card.NewCard(card.Five), DoubleDown},
		{"soft 17 vs 6", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Six)}, card.NewCard(card.Six), DoubleDown},
		{"soft 20 vs 10", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine)}, card.NewCard(card.Ten), Stand},
		{"pair 8s vs 6", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight)}, card.NewCard(card.Six), Split},
		{"pair 10s vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)}, card.NewCard(card.Ten), Stand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playerHand := hand.NewHand(tt.playerCards)
			dealerHand := hand.NewHand([]card.Card{tt.dealerCard})
			got, err := MakeDecision(playerHand, dealerHand)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.wantDecision {
				t.Errorf("MakeDecision() = %v, want %v", got, tt.wantDecision)
			}
		})
	}
}

func TestDealerCardIndex(t *testing.T) {
	tests := []struct {
		card card.Card
		want int
	}{
		{card.NewCard(card.Ace), 0},
		{card.NewCard(card.Ten), 1},
		{card.NewCard(card.Jack), 1},
		{card.NewCard(card.Queen), 1},
		{card.NewCard(card.King), 1},
		{card.NewCard(card.Nine), 2},
		{card.NewCard(card.Eight), 3},
		{card.NewCard(card.Seven), 4},
		{card.NewCard(card.Six), 5},
		{card.NewCard(card.Five), 6},
		{card.NewCard(card.Four), 7},
		{card.NewCard(card.Three), 8},
		{card.NewCard(card.Two), 9},
	}

	for _, tt := range tests {
		if got := dealerCardIndex(tt.card); got != tt.want {
			t.Errorf("dealerCardIndex(%v) = %d, want %d", tt.card, got, tt.want)
		}
	}
}
