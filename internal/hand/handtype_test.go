package hand

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
)

func TestHand_GetType_Hard(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  HandType
	}{
		{"hard 21", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Four)}, Hard21},
		{"hard 20", []card.Card{card.NewCard(card.King), card.NewCard(card.Queen)}, Hard20},
		{"hard 19", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Nine)}, Hard19},
		{"hard 18", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Eight)}, Hard18},
		{"hard 17", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}, Hard17},
		{"hard 16", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)}, Hard16},
		{"hard 15", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Five)}, Hard15},
		{"hard 14", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Four)}, Hard14},
		{"hard 13", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Three)}, Hard13},
		{"hard 12", []card.Card{card.NewCard(card.Seven), card.NewCard(card.Five)}, Hard12},
		{"hard 11", []card.Card{card.NewCard(card.Six), card.NewCard(card.Five)}, Hard11},
		{"hard 10", []card.Card{card.NewCard(card.Six), card.NewCard(card.Four)}, Hard10},
		{"hard 9", []card.Card{card.NewCard(card.Five), card.NewCard(card.Four)}, Hard9},
		{"hard 8", []card.Card{card.NewCard(card.Five), card.NewCard(card.Three)}, Hard8},
		{"hard 7", []card.Card{card.NewCard(card.Four), card.NewCard(card.Three)}, Hard7},
		{"hard 6", []card.Card{card.NewCard(card.Four), card.NewCard(card.Two)}, Hard6},
		{"hard 5", []card.Card{card.NewCard(card.Three), card.NewCard(card.Two)}, Hard5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.GetType(); got != tt.want {
				t.Errorf("Hand.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_GetType_Soft(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  HandType
	}{
		{"soft 20", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine)}, SoftA9},
		{"soft 19", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Eight)}, SoftA8},
		{"soft 18", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Seven)}, SoftA7},
		{"soft 17", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Six)}, SoftA6},
		{"soft 16", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Five)}, SoftA5},
		{"soft 15", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Four)}, SoftA4},
		{"soft 14", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Three)}, SoftA3},
		{"soft 13", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Two)}, SoftA2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.GetType(); got != tt.want {
				t.Errorf("Hand.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_GetType_Pair(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  HandType
	}{
		{"pair aces", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ace)}, PairA},
		{"pair 10s", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)}, Pair10},
		{"pair jacks", []card.Card{card.NewCard(card.Jack), card.NewCard(card.Jack)}, Pair10},
		{"pair queens", []card.Card{card.NewCard(card.Queen), card.NewCard(card.Queen)}, Pair10},
		{"pair kings", []card.Card{card.NewCard(card.King), card.NewCard(card.King)}, Pair10},
		{"pair 9s", []card.Card{card.NewCard(card.Nine), card.NewCard(card.Nine)}, Pair9},
		{"pair 8s", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight)}, Pair8},
		{"pair 7s", []card.Card{card.NewCard(card.Seven), card.NewCard(card.Seven)}, Pair7},
		{"pair 6s", []card.Card{card.NewCard(card.Six), card.NewCard(card.Six)}, Pair6},
		{"pair 5s", []card.Card{card.NewCard(card.Five), card.NewCard(card.Five)}, Pair5},
		{"pair 4s", []card.Card{card.NewCard(card.Four), card.NewCard(card.Four)}, Pair4},
		{"pair 3s", []card.Card{card.NewCard(card.Three), card.NewCard(card.Three)}, Pair3},
		{"pair 2s", []card.Card{card.NewCard(card.Two), card.NewCard(card.Two)}, Pair2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.GetType(); got != tt.want {
				t.Errorf("Hand.GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}
