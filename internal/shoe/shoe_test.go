package shoe

import (
	"math/rand"
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/deck"
)

func TestNewShoe(t *testing.T) {
	s := NewShoe([]deck.Deck{deck.NewDeck(), deck.NewDeck()})

	expectedCards := 2 * 52
	if len(s.Cards) != expectedCards {
		t.Errorf("NewShoe() Cards length = %d, want %d", len(s.Cards), expectedCards)
	}
}

func TestShoe_Shuffle(t *testing.T) {
	rng1 := rand.New(rand.NewSource(42))
	rng2 := rand.New(rand.NewSource(43))

	s1 := NewShoe([]deck.Deck{deck.NewDeck()})
	s2 := NewShoe([]deck.Deck{deck.NewDeck()})

	s1.Shuffle(rng1)
	s2.Shuffle(rng2)

	if len(s1.Cards) != 52 || len(s2.Cards) != 52 {
		t.Errorf("Shuffle() changed Cards length")
	}

	rankCounts1 := make(map[string]int)
	rankCounts2 := make(map[string]int)
	for _, c := range s1.Cards {
		rankCounts1[c.ToString()]++
	}
	for _, c := range s2.Cards {
		rankCounts2[c.ToString()]++
	}

	for rank, count := range rankCounts1 {
		if rankCounts2[rank] != count {
			t.Errorf("Shuffle() changed rank distribution for %s", rank)
		}
	}
}

func TestShoe_Shuffle_ChangesOrder(t *testing.T) {
	rng := rand.New(rand.NewSource(42))

	s1 := NewShoe([]deck.Deck{deck.NewDeck()})
	s2 := NewShoe([]deck.Deck{deck.NewDeck()})

	s2.Shuffle(rng)

	allSame := true
	cards1 := s1.Cards
	cards2 := s2.Cards
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

func TestShoe_Draw(t *testing.T) {
	s := NewShoe([]deck.Deck{deck.NewDeck()})
	initialTotal := len(s.Cards)

	drawn := s.Draw()

	if len(s.Cards) != initialTotal-1 {
		t.Errorf("Draw() Cards length = %d, want %d", len(s.Cards), initialTotal-1)
	}

	if drawn.Rank() != card.Two {
		t.Errorf("Draw() rank = %v, want %v", drawn.Rank(), card.Two)
	}
}

func TestShoe_Draw_Multiple(t *testing.T) {
	s := NewShoe([]deck.Deck{deck.NewDeck()})
	initialTotal := len(s.Cards)

	cardsDrawn := 10
	for i := 0; i < cardsDrawn; i++ {
		s.Draw()
	}

	if len(s.Cards) != initialTotal-cardsDrawn {
		t.Errorf("Draw() %d times: Cards length = %d, want %d", cardsDrawn, len(s.Cards), initialTotal-cardsDrawn)
	}
}
