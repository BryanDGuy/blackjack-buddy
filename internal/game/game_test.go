package game

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/shoe"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func cards(ranks ...card.Rank) []card.Card {
	cards := make([]card.Card, len(ranks))
	for i, rank := range ranks {
		cards[i] = card.NewCard(rank)
	}
	return cards
}

func testGame(player, dealer []card.Rank) *Game {
	g := NewGame()
	if player != nil {
		g.ActiveHand = hand.NewHand(cards(player...))
	}
	if dealer != nil {
		g.DealerHand = hand.NewHand(cards(dealer...))
	}
	return g
}

func TestNewGameStartsPlayableSixDeckRound(t *testing.T) {
	g := NewGame()
	if g.RoundState != RoundStateNone || len(g.shoe.Cards) != decksInShoe*52 {
		t.Fatalf("new game = %#v, shoe cards = %d", g, len(g.shoe.Cards))
	}
	g.StartRound()
	if g.ActiveHand == nil || g.DealerHand == nil || len(g.ActiveHand.Cards) != 2 || len(g.DealerHand.Cards) != 2 || len(g.shoe.Cards) != decksInShoe*52-4 {
		t.Fatalf("started game = %#v, shoe cards = %d", g, len(g.shoe.Cards))
	}
}

func TestRoundStateAndMoves(t *testing.T) {
	t.Run("abandon", func(t *testing.T) {
		g := NewGame()
		g.StartRound()
		g.AbandonRound()
		if g.RoundState != RoundStateComplete || g.ActiveHand != nil || len(g.ResolvedHands) != 0 {
			t.Fatalf("abandoned game = %#v", g)
		}
	})

	t.Run("no active hand", func(t *testing.T) {
		for _, move := range []strategy.Decision{strategy.Hit, strategy.Stand, strategy.DoubleDown, strategy.Split} {
			if err := NewGame().ApplyMove(move); !errors.Is(err, ErrNoActiveHand) {
				t.Errorf("%s error = %v, want %v", move, err, ErrNoActiveHand)
			}
		}
	})

	t.Run("hit", func(t *testing.T) {
		g := testGame([]card.Rank{card.Two, card.Three}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.Two)}
		if err := g.ApplyMove(strategy.Hit); err != nil || len(g.ActiveHand.Cards) != 3 {
			t.Fatalf("hit error = %v, hand = %#v", err, g.ActiveHand)
		}
	})

	t.Run("hit resolves on 21", func(t *testing.T) {
		g := testGame([]card.Rank{card.Ten, card.Nine}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.Two)}
		if err := g.ApplyMove(strategy.Hit); err != nil || g.ActiveHand != nil || g.RoundState != RoundStateComplete || !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeWin}) {
			t.Fatalf("hit result: err = %v, game = %#v", err, g)
		}
	})

	t.Run("hit resolves on bust", func(t *testing.T) {
		g := testGame([]card.Rank{card.Ten, card.Nine}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.King)}
		if err := g.ApplyMove(strategy.Hit); err != nil || g.ActiveHand != nil || !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeBust}) {
			t.Fatalf("hit result: err = %v, game = %#v", err, g)
		}
	})

	t.Run("stand completes dealer on S17", func(t *testing.T) {
		g := testGame([]card.Rank{card.Ten, card.Seven}, []card.Rank{card.Ten, card.Six})
		g.shoe = shoe.Shoe{Cards: cards(card.Two)}
		if err := g.ApplyMove(strategy.Stand); err != nil || g.DealerHand.Value() != 18 || !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeLose}) {
			t.Fatalf("stand result: err = %v, game = %#v", err, g)
		}
	})

	t.Run("double only with two cards", func(t *testing.T) {
		g := testGame([]card.Rank{card.Three, card.Four}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.Two)}
		if err := g.ApplyMove(strategy.Hit); err != nil {
			t.Fatal(err)
		}
		if err := g.ApplyMove(strategy.DoubleDown); !errors.Is(err, ErrInvalidMove) {
			t.Fatalf("double error = %v, want %v", err, ErrInvalidMove)
		}
	})

	t.Run("double completes the hand", func(t *testing.T) {
		g := testGame([]card.Rank{card.Ten, card.Seven}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.Two)}
		if err := g.ApplyMove(strategy.DoubleDown); err != nil || len(g.ResolvedHands) != 1 || len(g.ResolvedHands[0].Cards) != 3 || !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeWin}) {
			t.Fatalf("double result: err = %v, game = %#v", err, g)
		}
	})

	t.Run("split advances and permits double", func(t *testing.T) {
		g := testGame([]card.Rank{card.Eight, card.Eight}, []card.Rank{card.Ten, card.Seven})
		g.shoe = shoe.Shoe{Cards: cards(card.Three, card.Four, card.Two)}
		if err := g.ApplyMove(strategy.Split); err != nil {
			t.Fatal(err)
		}
		if len(g.UnresolvedHands) != 1 || !g.ActiveHand.FromSplit || !g.UnresolvedHands[0].FromSplit {
			t.Fatalf("split game = %#v", g)
		}
		second := g.UnresolvedHands[0]
		if err := g.ApplyMove(strategy.DoubleDown); err != nil || len(g.ResolvedHands) != 1 || g.ActiveHand != second || len(g.UnresolvedHands) != 0 || len(g.DealerHand.Cards) != 2 {
			t.Fatalf("split double: err = %v, game = %#v", err, g)
		}
	})

	t.Run("invalid split", func(t *testing.T) {
		for _, player := range [][]card.Rank{nil, {card.Ten, card.Seven}, {card.Ten, card.Ten, card.Two}} {
			g := testGame(player, []card.Rank{card.Ten, card.Seven})
			want := ErrInvalidSplit
			if player == nil {
				want = ErrNoActiveHand
			}
			if err := g.ApplyMove(strategy.Split); !errors.Is(err, want) {
				t.Errorf("split %v error = %v, want %v", player, err, want)
			}
		}
	})
}

func TestSplitTwentyOneIsNotBlackjack(t *testing.T) {
	g := testGame([]card.Rank{card.Ace, card.Ace}, []card.Rank{card.Ten, card.Nine})
	g.shoe = shoe.Shoe{Cards: cards(card.Ten, card.Ten)}
	for _, move := range []strategy.Decision{strategy.Split, strategy.Stand, strategy.Stand} {
		if err := g.ApplyMove(move); err != nil {
			t.Fatal(err)
		}
	}
	if !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeWin, OutcomeWin}) {
		t.Fatalf("outcomes = %v", g.Outcomes)
	}
}

func TestDealerCompletion(t *testing.T) {
	g := testGame(nil, []card.Rank{card.Ten, card.Six})
	g.shoe = shoe.Shoe{Cards: cards(card.Two)}
	g.completeDealerHand()
	if g.DealerHand.Value() != 18 {
		t.Fatalf("dealer value = %d, want 18", g.DealerHand.Value())
	}

	g = testGame(nil, []card.Rank{card.Ace, card.Six})
	g.shoe = shoe.Shoe{Cards: cards(card.Two)}
	g.completeDealerHand()
	if len(g.DealerHand.Cards) != 2 {
		t.Fatalf("soft 17 dealer cards = %d, want 2", len(g.DealerHand.Cards))
	}

	g = NewGame()
	g.completeDealerHand()
}

func TestOutcomes(t *testing.T) {
	tests := []struct {
		name   string
		player []card.Rank
		dealer []card.Rank
		split  bool
		want   Outcome
	}{
		{"bust", []card.Rank{card.Ten, card.Seven, card.King}, []card.Rank{card.Ten, card.Eight}, false, OutcomeBust},
		{"win", []card.Rank{card.Ten, card.Eight}, []card.Rank{card.Ten, card.Seven}, false, OutcomeWin},
		{"lose", []card.Rank{card.Ten, card.Seven}, []card.Rank{card.Ten, card.Eight}, false, OutcomeLose},
		{"push", []card.Rank{card.Ten, card.Seven}, []card.Rank{card.Ten, card.Seven}, false, OutcomePush},
		{"blackjack", []card.Rank{card.Ten, card.Ace}, []card.Rank{card.Ten, card.Seven}, false, OutcomeBlackjack},
		{"blackjack push", []card.Rank{card.Ten, card.Ace}, []card.Rank{card.Ten, card.Ace}, false, OutcomePush},
		{"dealer bust", []card.Rank{card.Ten, card.Seven}, []card.Rank{card.Ten, card.Seven, card.King}, false, OutcomeWin},
		{"split 21", []card.Rank{card.Ace, card.Ten}, []card.Rank{card.Ten, card.Nine}, true, OutcomeWin},
		{"dealer blackjack beats split 21", []card.Rank{card.Ace, card.Ten}, []card.Rank{card.Ten, card.Ace}, true, OutcomeLose},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := testGame(nil, tt.dealer)
			h := hand.NewHand(cards(tt.player...))
			h.FromSplit = tt.split
			g.ResolvedHands = []*hand.Hand{h}
			g.setOutcomes()
			if !reflect.DeepEqual(g.Outcomes, []Outcome{tt.want}) {
				t.Fatalf("outcomes = %v, want %v", g.Outcomes, tt.want)
			}
		})
	}

	g := testGame(nil, nil)
	g.ResolvedHands = []*hand.Hand{hand.NewHand(cards(card.Ten, card.Seven))}
	g.setOutcomes()
	if len(g.Outcomes) != 0 {
		t.Fatalf("nil dealer outcomes = %v", g.Outcomes)
	}

	g = testGame(nil, []card.Rank{card.Ten, card.Eight})
	g.ResolvedHands = []*hand.Hand{hand.NewHand(cards(card.Ten, card.Seven)), hand.NewHand(cards(card.Ten, card.King))}
	g.setOutcomes()
	if !reflect.DeepEqual(g.Outcomes, []Outcome{OutcomeLose, OutcomeWin}) {
		t.Fatalf("multiple outcomes = %v", g.Outcomes)
	}
}
