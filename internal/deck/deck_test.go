package deck

import (
	"math/rand"
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

func TestDeck_Shuffle(t *testing.T) {
	rng1 := rand.New(rand.NewSource(42))
	rng2 := rand.New(rand.NewSource(43))

	d1 := NewDeck()
	d2 := NewDeck()

	d1.Shuffle(rng1)
	d2.Shuffle(rng2)

	if len(d1.Cards) != 52 || len(d2.Cards) != 52 {
		t.Errorf("Shuffle() changed deck size")
	}

	rankCounts1 := make(map[card.Rank]int)
	rankCounts2 := make(map[card.Rank]int)

	for _, c := range d1.Cards {
		rankCounts1[c.Rank()]++
	}
	for _, c := range d2.Cards {
		rankCounts2[c.Rank()]++
	}

	for rank, count := range rankCounts1 {
		if rankCounts2[rank] != count {
			t.Errorf("Shuffle() changed rank distribution for %v", rank)
		}
	}
}

func TestDeck_Shuffle_ChangesOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	d1 := NewDeck()
	d2 := NewDeck()

	d2.Shuffle(rng)

	allSame := true
	cards1 := d1.Cards
	cards2 := d2.Cards
	for i := range cards1 {
		if cards1[i].Rank() != cards2[i].Rank() {
			allSame = false
			break
		}
	}

	if allSame {
		t.Errorf("Shuffle() did not change card order")
	}
}
