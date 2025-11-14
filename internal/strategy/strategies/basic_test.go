package strategies

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func TestBasic_GetDecision(t *testing.T) {
	s := NewBasic()

	tests := []struct {
		name         string
		playerCards  []card.Card
		dealerCard   card.Card
		wantDecision strategy.Decision
	}{
		{"hard 20 vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)}, card.NewCard(card.Ten), strategy.Stand},
		{"hard 16 vs 6", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)}, card.NewCard(card.Six), strategy.Stand},
		{"hard 16 vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)}, card.NewCard(card.Ten), strategy.Hit},
		{"hard 11 vs 5", []card.Card{card.NewCard(card.Seven), card.NewCard(card.Four)}, card.NewCard(card.Five), strategy.DoubleDown},
		{"soft 17 vs 6", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Six)}, card.NewCard(card.Six), strategy.DoubleDown},
		{"soft 20 vs 10", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine)}, card.NewCard(card.Ten), strategy.Stand},
		{"pair 8s vs 6", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight)}, card.NewCard(card.Six), strategy.Split},
		{"pair 10s vs 10", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)}, card.NewCard(card.Ten), strategy.Stand},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playerHand := hand.NewHand(tt.playerCards)
			dealerHand := hand.NewHand([]card.Card{tt.dealerCard})
			if got := s.GetDecision(playerHand, dealerHand); got != tt.wantDecision {
				t.Errorf("Basic.GetDecision() = %v, want %v", got, tt.wantDecision)
			}
		})
	}
}

func TestGetDealerCardIndex(t *testing.T) {
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
		if got := GetDealerCardIndex(tt.card); got != tt.want {
			t.Errorf("GetDealerCardIndex(%v) = %d, want %d", tt.card, got, tt.want)
		}
	}
}
