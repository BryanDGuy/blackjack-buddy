package game

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/player"
)

func TestGame_Hit(t *testing.T) {
	t.Run("hit with empty hand", func(t *testing.T) {
		p := player.NewPlayer()
		g := NewGame(p, dealer.NewDealer())

		err := g.Hit()

		if err != player.ErrNoActiveHand {
			t.Errorf("Hit() error = %v, want %v", err, player.ErrNoActiveHand)
		}
	})

	t.Run("hit adds card", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Two), card.NewCard(card.Three)})
		g := NewGame(p, dealer.NewDealer())

		err := g.Hit()

		if err != nil {
			t.Errorf("Hit() error = %v, want nil", err)
		}

		if g.Player.ActiveHand != nil {
			if len(g.Player.ActiveHand.Cards) != 3 {
				t.Errorf("Hit() should add one card, got %d cards, want 3", len(g.Player.ActiveHand.Cards))
			}
		} else {
			if len(g.Player.ResolvedHands) == 0 {
				t.Error("Hit() should complete hand if bust or 21, but no resolved hands found")
			}
		}
	})
}

func TestGame_Stand(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

	g := NewGame(p, dealer.NewDealer())
	err := g.Stand()

	if err != nil {
		t.Errorf("Stand() error = %v, want nil", err)
	}

	if g.Player.ActiveHand != nil {
		t.Error("Stand() should complete hand")
	}

	if len(g.Player.ResolvedHands) != 1 {
		t.Errorf("Stand() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
	}
}

func TestGame_Stand_NoActiveHand(t *testing.T) {
	p := player.NewPlayer()
	g := NewGame(p, dealer.NewDealer())

	err := g.Stand()

	if err != player.ErrNoActiveHand {
		t.Errorf("Stand() error = %v, want %v", err, player.ErrNoActiveHand)
	}
}

func TestGame_Double(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

	g := NewGame(p, dealer.NewDealer())
	err := g.Double()

	if err != nil {
		t.Errorf("Double() error = %v, want nil", err)
	}

	if g.Player.ActiveHand != nil {
		t.Error("Double() should complete hand")
	}

	if len(g.Player.ResolvedHands) != 1 {
		t.Errorf("Double() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
	}

	if len(g.Player.ResolvedHands[0].Cards) != 3 {
		t.Errorf("Double() should have 3 cards in resolved hand, got %d", len(g.Player.ResolvedHands[0].Cards))
	}
}

func TestGame_Double_InvalidMove(t *testing.T) {
	tests := []struct {
		name        string
		initialHand []card.Card
	}{
		{
			name:        "no active hand",
			initialHand: nil,
		},
		{
			name:        "more than 2 cards",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two)},
		},
		{
			name:        "single card",
			initialHand: []card.Card{card.NewCard(card.Ten)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := player.NewPlayer()
			if tt.initialHand != nil {
				p.ActiveHand = hand.NewHand(tt.initialHand)
			}

			g := NewGame(p, dealer.NewDealer())
			err := g.Double()

			if err != player.ErrInvalidMove {
				t.Errorf("Double() error = %v, want %v", err, player.ErrInvalidMove)
			}
		})
	}
}

func TestGame_Split(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)})

	g := NewGame(p, dealer.NewDealer())
	err := g.Split()

	if err != nil {
		t.Errorf("Split() error = %v, want nil", err)
	}

	if g.Player.ActiveHand == nil {
		t.Fatal("Split() should have active hand")
	}

	if len(g.Player.ActiveHand.Cards) != 2 {
		t.Errorf("Split() active hand should have 2 cards, got %d", len(g.Player.ActiveHand.Cards))
	}

	if len(g.Player.UnresolvedHands) != 1 {
		t.Errorf("Split() should have 1 unresolved hand, got %d", len(g.Player.UnresolvedHands))
	}

	if len(g.Player.UnresolvedHands[0].Cards) != 2 {
		t.Errorf("Split() unresolved hand should have 2 cards, got %d", len(g.Player.UnresolvedHands[0].Cards))
	}
}

func TestGame_Split_InvalidMove(t *testing.T) {
	tests := []struct {
		name        string
		initialHand []card.Card
		wantErr     error
	}{
		{
			name:        "no active hand",
			initialHand: nil,
			wantErr:     player.ErrInvalidMove,
		},
		{
			name:        "more than 2 cards",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two)},
			wantErr:     player.ErrInvalidMove,
		},
		{
			name:        "cannot split different ranks",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			wantErr:     player.ErrInvalidSplit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := player.NewPlayer()
			if tt.initialHand != nil {
				p.ActiveHand = hand.NewHand(tt.initialHand)
			}

			g := NewGame(p, dealer.NewDealer())
			err := g.Split()

			if err != tt.wantErr {
				t.Errorf("Split() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_completeAndAdvance(t *testing.T) {
	t.Run("advance to unresolved hand", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})
		p.UnresolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Nine), card.NewCard(card.Eight)}),
		}

		g := NewGame(p, dealer.NewDealer())
		g.completeAndAdvance()

		if g.Player.ActiveHand == nil {
			t.Fatal("completeAndAdvance() should have active hand after advancing")
		}

		if len(g.Player.ResolvedHands) != 1 {
			t.Errorf("completeAndAdvance() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
		}

		if len(g.Player.UnresolvedHands) != 0 {
			t.Errorf("completeAndAdvance() should have 0 unresolved hands, got %d", len(g.Player.UnresolvedHands))
		}
	})

	t.Run("no unresolved hands", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		g := NewGame(p, dealer.NewDealer())
		g.completeAndAdvance()

		if g.Player.ActiveHand != nil {
			t.Error("completeAndAdvance() should not have active hand when no unresolved hands")
		}

		if len(g.Player.ResolvedHands) != 1 {
			t.Errorf("completeAndAdvance() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
		}
	})
}
