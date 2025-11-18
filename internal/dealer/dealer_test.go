package dealer

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type mockCardDrawer struct {
	cards []card.Card
	index int
}

func (m *mockCardDrawer) DrawCard() card.Card {
	if m.index >= len(m.cards) {
		return card.NewCard(card.Two)
	}
	card := m.cards[m.index]
	m.index++
	return card
}

func TestNewDealer(t *testing.T) {
	d := NewDealer()
	if d == nil {
		t.Fatal("NewDealer returned nil")
	}
	if d.Hand != nil {
		t.Error("NewDealer should have nil Hand")
	}
}

func TestFinish_WithOneCard(t *testing.T) {
	d := NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten)})

	drawer := &mockCardDrawer{
		cards: []card.Card{card.NewCard(card.Seven)},
	}

	d.Finish(drawer)

	if len(d.Hand.Cards) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(d.Hand.Cards))
	}
}

func TestFinish_StopsAt17(t *testing.T) {
	d := NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Six)})

	drawer := &mockCardDrawer{
		cards: []card.Card{card.NewCard(card.Ace)},
	}

	d.Finish(drawer)

	if d.Hand.Value() < 17 {
		t.Errorf("Expected hand value >= 17, got %d", d.Hand.Value())
	}
}

func TestFinish_ContinuesUntil17(t *testing.T) {
	d := NewDealer()
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Two), card.NewCard(card.Three)})

	drawer := &mockCardDrawer{
		cards: []card.Card{
			card.NewCard(card.Four),
			card.NewCard(card.Five),
			card.NewCard(card.Six),
		},
	}

	d.Finish(drawer)

	if d.Hand.Value() < 17 {
		t.Errorf("Expected hand value >= 17, got %d", d.Hand.Value())
	}
}

func TestFinish_WithNilHand(t *testing.T) {
	d := NewDealer()
	d.Hand = nil

	drawer := &mockCardDrawer{cards: []card.Card{card.NewCard(card.Ten)}}

	d.Finish(drawer)

	if d.Hand != nil {
		t.Error("Finish should not modify nil Hand")
	}
}

func TestFinish_WithEmptyHand(t *testing.T) {
	d := NewDealer()
	d.Hand = hand.NewHand(nil)

	drawer := &mockCardDrawer{cards: []card.Card{card.NewCard(card.Ten)}}

	d.Finish(drawer)

	if len(d.Hand.Cards) != 0 {
		t.Error("Finish should not modify empty Hand")
	}
}
