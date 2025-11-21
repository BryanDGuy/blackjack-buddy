package player

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()
	if p == nil {
		t.Fatal("NewPlayer() returned nil")
	}
	if p.ActiveHand != nil {
		t.Error("NewPlayer() ActiveHand should be nil")
	}
	if p.UnresolvedHands != nil {
		t.Error("NewPlayer() UnresolvedHands should be nil")
	}
	if p.ResolvedHands != nil {
		t.Error("NewPlayer() ResolvedHands should be nil")
	}
}

func TestPlayer_CanMove(t *testing.T) {
	tests := []struct {
		name        string
		activeHand  []card.Card
		wantCanMove bool
	}{
		{
			name:        "no active hand",
			activeHand:  nil,
			wantCanMove: false,
		},
		{
			name:        "can move with valid hand",
			activeHand:  []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			wantCanMove: true,
		},
		{
			name:        "cannot move with bust",
			activeHand:  []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.King)},
			wantCanMove: false,
		},
		{
			name:        "cannot move with empty hand",
			activeHand:  []card.Card{},
			wantCanMove: false,
		},
		{
			name:        "cannot move with 21",
			activeHand:  []card.Card{card.NewCard(card.Ten), card.NewCard(card.Ace)},
			wantCanMove: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPlayer()
			if tt.activeHand != nil {
				p.ActiveHand = hand.NewHand(tt.activeHand)
			}

			got := p.CanMove()

			if got != tt.wantCanMove {
				t.Errorf("CanMove() = %v, want %v", got, tt.wantCanMove)
			}
		})
	}
}
