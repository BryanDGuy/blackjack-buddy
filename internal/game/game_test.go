package game

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/player"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func TestGame_ApplyMove_Hit(t *testing.T) {
	t.Run("hit with empty hand", func(t *testing.T) {
		p := player.NewPlayer()
		g := NewGame(p, dealer.NewDealer())

		err := g.ApplyMove(strategy.Hit)

		if err != ErrNoActiveHand {
			t.Errorf("Hit() error = %v, want %v", err, ErrNoActiveHand)
		}
	})

	t.Run("hit adds card", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Two), card.NewCard(card.Three)})
		g := NewGame(p, dealer.NewDealer())

		err := g.ApplyMove(strategy.Hit)

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

func TestGame_ApplyMove_Stand(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

	d := dealer.NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

	g := NewGame(p, d)
	err := g.ApplyMove(strategy.Stand)

	if err != nil {
		t.Errorf("Stand() error = %v, want nil", err)
	}

	if g.Player.ActiveHand != nil {
		t.Error("Stand() should complete hand")
	}

	if len(g.Player.ResolvedHands) != 1 {
		t.Errorf("Stand() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
	}

	if d.Hand.Value() < 17 {
		t.Errorf("Stand() should complete dealer hand (value >= 17), got %d", d.Hand.Value())
	}

	if len(g.Outcomes) != 1 {
		t.Errorf("Stand() should set outcomes, got %d outcomes", len(g.Outcomes))
	}
}

func TestGame_ApplyMove_Stand_NoActiveHand(t *testing.T) {
	p := player.NewPlayer()
	g := NewGame(p, dealer.NewDealer())

	err := g.ApplyMove(strategy.Stand)

	if err != ErrNoActiveHand {
		t.Errorf("Stand() error = %v, want %v", err, ErrNoActiveHand)
	}
}

func TestGame_ApplyMove_Double(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

	d := dealer.NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

	g := NewGame(p, d)
	err := g.ApplyMove(strategy.DoubleDown)

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

	if d.Hand.Value() < 17 {
		t.Errorf("Double() should complete dealer hand (value >= 17), got %d", d.Hand.Value())
	}

	if len(g.Outcomes) != 1 {
		t.Errorf("Double() should set outcomes, got %d outcomes", len(g.Outcomes))
	}
}

func TestGame_ApplyMove_Double_InvalidMove(t *testing.T) {
	tests := []struct {
		name        string
		initialHand []card.Card
		wantErr     error
	}{
		{
			name:        "no active hand",
			initialHand: nil,
			wantErr:     ErrNoActiveHand,
		},
		{
			name:        "more than 2 cards",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two)},
			wantErr:     ErrInvalidMove,
		},
		{
			name:        "single card",
			initialHand: []card.Card{card.NewCard(card.Ten)},
			wantErr:     ErrInvalidMove,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := player.NewPlayer()
			if tt.initialHand != nil {
				p.ActiveHand = hand.NewHand(tt.initialHand)
			}

			g := NewGame(p, dealer.NewDealer())
			err := g.ApplyMove(strategy.DoubleDown)

			if err != tt.wantErr {
				t.Errorf("Double() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ApplyMove_Split(t *testing.T) {
	p := player.NewPlayer()
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten)})

	g := NewGame(p, dealer.NewDealer())
	err := g.ApplyMove(strategy.Split)

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

func TestGame_ApplyMove_Split_InvalidMove(t *testing.T) {
	tests := []struct {
		name        string
		initialHand []card.Card
		wantErr     error
	}{
		{
			name:        "no active hand",
			initialHand: nil,
			wantErr:     ErrNoActiveHand,
		},
		{
			name:        "more than 2 cards",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two)},
			wantErr:     ErrInvalidMove,
		},
		{
			name:        "cannot split different ranks",
			initialHand: []card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			wantErr:     ErrInvalidSplit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := player.NewPlayer()
			if tt.initialHand != nil {
				p.ActiveHand = hand.NewHand(tt.initialHand)
			}

			g := NewGame(p, dealer.NewDealer())
			err := g.ApplyMove(strategy.Split)

			if err != tt.wantErr {
				t.Errorf("Split() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestGame_ApplyMove_Advancement(t *testing.T) {
	t.Run("advance to unresolved hand does not complete dealer", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})
		p.UnresolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Nine), card.NewCard(card.Eight)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

		g := NewGame(p, d)
		err := g.ApplyMove(strategy.Stand)

		if err != nil {
			t.Fatalf("ApplyMove() error = %v, want nil", err)
		}

		if g.Player.ActiveHand == nil {
			t.Fatal("ApplyMove() should have active hand after advancing")
		}

		if len(g.Player.ResolvedHands) != 1 {
			t.Errorf("ApplyMove() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
		}

		if len(g.Player.UnresolvedHands) != 0 {
			t.Errorf("ApplyMove() should have 0 unresolved hands, got %d", len(g.Player.UnresolvedHands))
		}

		if len(d.Hand.Cards) != 2 {
			t.Errorf("ApplyMove() should not complete dealer hand when unresolved hands exist, dealer should still have 2 cards, got %d", len(d.Hand.Cards))
		}

		if len(g.Outcomes) != 0 {
			t.Errorf("ApplyMove() should not set outcomes when unresolved hands exist, got %d outcomes", len(g.Outcomes))
		}
	})

	t.Run("no unresolved hands completes dealer and sets outcomes", func(t *testing.T) {
		p := player.NewPlayer()
		p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

		g := NewGame(p, d)
		err := g.ApplyMove(strategy.Stand)

		if err != nil {
			t.Fatalf("ApplyMove() error = %v, want nil", err)
		}

		if g.Player.ActiveHand != nil {
			t.Error("ApplyMove() should not have active hand when no unresolved hands")
		}

		if len(g.Player.ResolvedHands) != 1 {
			t.Errorf("ApplyMove() should have 1 resolved hand, got %d", len(g.Player.ResolvedHands))
		}

		if d.Hand.Value() < 17 {
			t.Errorf("ApplyMove() should complete dealer hand (value >= 17), got %d", d.Hand.Value())
		}

		if len(g.Outcomes) != 1 {
			t.Errorf("ApplyMove() should set outcomes when player is done, got %d outcomes", len(g.Outcomes))
		}
	})
}

func TestGame_completeDealerHand(t *testing.T) {
	t.Run("dealer hits until 17 or above", func(t *testing.T) {
		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

		g := NewGame(player.NewPlayer(), d)
		g.completeDealerHand()

		if d.Hand.Value() < 17 {
			t.Errorf("completeDealerHand() should hit until value >= 17, got %d", d.Hand.Value())
		}
	})

	t.Run("dealer stands on 17", func(t *testing.T) {
		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		initialCardCount := len(d.Hand.Cards)

		g := NewGame(player.NewPlayer(), d)
		g.completeDealerHand()

		if len(d.Hand.Cards) != initialCardCount {
			t.Errorf("completeDealerHand() should not hit when value is 17, got %d cards", len(d.Hand.Cards))
		}

		if d.Hand.Value() != 17 {
			t.Errorf("completeDealerHand() should keep value at 17, got %d", d.Hand.Value())
		}
	})

	t.Run("dealer busts", func(t *testing.T) {
		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

		g := NewGame(player.NewPlayer(), d)

		for i := 0; i < 10 && d.Hand.Value() < 17; i++ {
			g.completeDealerHand()
		}

		if d.Hand.Value() < 17 {
			t.Errorf("completeDealerHand() should continue until value >= 17 or bust, but value is %d", d.Hand.Value())
		}
	})

	t.Run("no dealer hand", func(t *testing.T) {
		d := dealer.NewDealer()

		g := NewGame(player.NewPlayer(), d)
		g.completeDealerHand()

		if d.Hand != nil {
			t.Error("completeDealerHand() should handle nil dealer hand gracefully")
		}
	})
}

func TestGame_setOutcomes(t *testing.T) {
	t.Run("bust hand", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.King)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Eight)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeBust {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomeBust)
		}
	})

	t.Run("player wins", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Eight)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeWin {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomeWin)
		}
	})

	t.Run("player loses", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Eight)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeLose {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomeLose)
		}
	})

	t.Run("push", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomePush {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomePush)
		}
	})

	t.Run("blackjack", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ace)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeBlackjack {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomeBlackjack)
		}
	})

	t.Run("blackjack push", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ace)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ace)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomePush {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomePush)
		}
	})

	t.Run("dealer bust", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.King)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 1 {
			t.Fatalf("SetOutcomes() should have 1 outcome, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeWin {
			t.Errorf("SetOutcomes() outcome = %v, want %v", g.Outcomes[0], OutcomeWin)
		}
	})

	t.Run("no dealer hand", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}),
		}

		g := NewGame(p, dealer.NewDealer())
		g.setOutcomes()

		if len(g.Outcomes) != 0 {
			t.Errorf("SetOutcomes() should have 0 outcomes when no dealer hand, got %d", len(g.Outcomes))
		}
	})

	t.Run("multiple hands", func(t *testing.T) {
		p := player.NewPlayer()
		p.ResolvedHands = []*hand.Hand{
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)}),
			hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.King)}),
		}

		d := dealer.NewDealer()
		d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Eight)})

		g := NewGame(p, d)
		g.setOutcomes()

		if len(g.Outcomes) != 2 {
			t.Fatalf("SetOutcomes() should have 2 outcomes, got %d", len(g.Outcomes))
		}

		if g.Outcomes[0] != OutcomeLose {
			t.Errorf("SetOutcomes() first outcome = %v, want %v", g.Outcomes[0], OutcomeLose)
		}

		if g.Outcomes[1] != OutcomeWin {
			t.Errorf("SetOutcomes() second outcome = %v, want %v", g.Outcomes[1], OutcomeWin)
		}
	})
}
