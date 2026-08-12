package deck

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
)

func TestNewDeck(t *testing.T) {
	d := NewDeck()

	if len(d.Cards) != 52 {
		t.Errorf("NewDeck() created deck with %d cards, want 52", len(d.Cards))
	}

	rankCounts := make(map[card.Rank]int)
	for _, c := range d.Cards {
		rankCounts[c.Rank()]++
	}

	expectedRanks := []card.Rank{
		card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King, card.Ace,
	}

	for _, rank := range expectedRanks {
		if count := rankCounts[rank]; count != 4 {
			t.Errorf("NewDeck() has %d cards of rank %v, want 4", count, rank)
		}
	}

	if len(rankCounts) != 13 {
		t.Errorf("NewDeck() has %d unique ranks, want 13", len(rankCounts))
	}
}
