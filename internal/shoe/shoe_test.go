package shoe

import (
	"math/rand"
	"testing"

	"github.com/bryan/blackjack-buddy/internal/deck"
)

func TestNewShoe(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	decks := make([]deck.Deck, 0, 2)
	for range 2 {
		d := deck.NewDeck(rng)
		d.Shuffle(rng)
		decks = append(decks, d)
	}

	s := NewShoe(decks)

	expectedCards := 2 * 52
	if len(s.Cards) != expectedCards {
		t.Errorf("NewShoe() Cards length = %d, want %d", len(s.Cards), expectedCards)
	}
}

func TestShoe_Shuffle(t *testing.T) {
	rng1 := rand.New(rand.NewSource(42))
	rng2 := rand.New(rand.NewSource(43))

	decks1 := make([]deck.Deck, 0, 1)
	decks2 := make([]deck.Deck, 0, 1)
	for range 1 {
		d1 := deck.NewDeck(rng1)
		d1.Shuffle(rng1)
		decks1 = append(decks1, d1)

		d2 := deck.NewDeck(rng2)
		d2.Shuffle(rng2)
		decks2 = append(decks2, d2)
	}

	s1 := NewShoe(decks1)
	s2 := NewShoe(decks2)

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

	decks := make([]deck.Deck, 0, 1)
	d := deck.NewDeck(rng)
	d.Shuffle(rng)
	decks = append(decks, d)

	s1 := NewShoe(decks)
	s2 := NewShoe(decks)

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
	rng := rand.New(rand.NewSource(42))
	decks := make([]deck.Deck, 0, 1)
	d := deck.NewDeck(rng)
	d.Shuffle(rng)
	decks = append(decks, d)

	s := NewShoe(decks)
	initialTotal := len(s.Cards)

	card := s.Draw()

	if len(s.Cards) != initialTotal-1 {
		t.Errorf("Draw() Cards length = %d, want %d", len(s.Cards), initialTotal-1)
	}

	if card.Rank() == 0 {
		t.Errorf("Draw() returned invalid card")
	}
}

func TestShoe_Draw_Multiple(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	decks := make([]deck.Deck, 0, 1)
	d := deck.NewDeck(rng)
	d.Shuffle(rng)
	decks = append(decks, d)

	s := NewShoe(decks)
	initialTotal := len(s.Cards)

	cardsDrawn := 10
	for i := 0; i < cardsDrawn; i++ {
		s.Draw()
	}

	if len(s.Cards) != initialTotal-cardsDrawn {
		t.Errorf("Draw() %d times: Cards length = %d, want %d", cardsDrawn, len(s.Cards), initialTotal-cardsDrawn)
	}
}
