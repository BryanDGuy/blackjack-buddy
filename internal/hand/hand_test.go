package hand

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
)

func TestHand_Value(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  int
	}{
		{"empty", []card.Card{}, 0},
		{"single card", []card.Card{card.NewCard(card.Five)}, 5},
		{"two cards", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}, 17},
		{"ace as 11", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Six)}, 17},
		{"ace as 1", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten), card.NewCard(card.Five)}, 16},
		{"multiple aces", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ace), card.NewCard(card.Nine)}, 21},
		{"bust with aces", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ace), card.NewCard(card.Ten), card.NewCard(card.Ten)}, 22},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.Value(); got != tt.want {
				t.Errorf("Hand.Value() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestHand_IsBust(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  bool
	}{
		{"not bust", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}, false},
		{"bust", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Five)}, true},
		{"exactly 21", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Four)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.IsBust(); got != tt.want {
				t.Errorf("Hand.IsBust() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_IsBlackjack(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  bool
	}{
		{"blackjack ace+ten", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten)}, true},
		{"blackjack ace+jack", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Jack)}, true},
		{"not blackjack", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine)}, false},
		{"three cards 21", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine), card.NewCard(card.Ace)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.IsBlackjack(); got != tt.want {
				t.Errorf("Hand.IsBlackjack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_CanSplit(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  bool
	}{
		{"can split", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight)}, true},
		{"cannot split", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Nine)}, false},
		{"single card", []card.Card{card.NewCard(card.Eight)}, false},
		{"three cards", []card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight), card.NewCard(card.Eight)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.CanSplit(); got != tt.want {
				t.Errorf("Hand.CanSplit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHand_IsSoft(t *testing.T) {
	tests := []struct {
		name  string
		cards []card.Card
		want  bool
	}{
		{"soft 17", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Six)}, true},
		{"soft 20", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Nine)}, true},
		{"hard 17", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}, false},
		{"blackjack", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten)}, false},
		{"bust", []card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten), card.NewCard(card.Five)}, false},
		{"no aces", []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHand(tt.cards)
			if got := h.IsSoft(); got != tt.want {
				t.Errorf("Hand.IsSoft() = %v, want %v", got, tt.want)
			}
		})
	}
}
