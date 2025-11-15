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
	if s.TotalCards != expectedCards {
		t.Errorf("NewShoe() TotalCards = %d, want %d", s.TotalCards, expectedCards)
	}

	if len(s.Cards) != expectedCards {
		t.Errorf("NewShoe() Cards length = %d, want %d", len(s.Cards), expectedCards)
	}

	rankCounts := make(map[string]int)
	for _, c := range s.Cards {
		rankStr := c.ToString()
		rankCounts[rankStr]++
	}

	for rank, count := range rankCounts {
		if s.RankCounts[rank] != count {
			t.Errorf("NewShoe() RankCounts[%s] = %d, want %d", rank, s.RankCounts[rank], count)
		}
	}

	for rank, count := range s.RankCounts {
		if rankCounts[rank] != count {
			t.Errorf("NewShoe() RankCounts[%s] = %d, but actual count is %d", rank, count, rankCounts[rank])
		}
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

	if s1.TotalCards != 52 || s2.TotalCards != 52 {
		t.Errorf("Shuffle() changed TotalCards")
	}

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
	for i := range s1.Cards {
		if s1.Cards[i].Rank != s2.Cards[i].Rank {
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
	initialTotal := s.TotalCards
	initialRankCounts := make(map[string]int)
	for k, v := range s.RankCounts {
		initialRankCounts[k] = v
	}

	card := s.Draw()

	if s.TotalCards != initialTotal-1 {
		t.Errorf("Draw() TotalCards = %d, want %d", s.TotalCards, initialTotal-1)
	}

	if len(s.Cards) != initialTotal-1 {
		t.Errorf("Draw() Cards length = %d, want %d", len(s.Cards), initialTotal-1)
	}

	rankStr := card.ToString()
	expectedCount := initialRankCounts[rankStr] - 1
	if s.RankCounts[rankStr] != expectedCount {
		t.Errorf("Draw() RankCounts[%s] = %d, want %d", rankStr, s.RankCounts[rankStr], expectedCount)
	}

	for rank, count := range initialRankCounts {
		if rank != rankStr {
			if s.RankCounts[rank] != count {
				t.Errorf("Draw() RankCounts[%s] = %d, want %d (should be unchanged)", rank, s.RankCounts[rank], count)
			}
		}
	}
}

func TestShoe_Draw_Multiple(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	decks := make([]deck.Deck, 0, 1)
	d := deck.NewDeck(rng)
	d.Shuffle(rng)
	decks = append(decks, d)

	s := NewShoe(decks)
	initialTotal := s.TotalCards

	cardsDrawn := 10
	for i := 0; i < cardsDrawn; i++ {
		s.Draw()
	}

	if s.TotalCards != initialTotal-cardsDrawn {
		t.Errorf("Draw() %d times: TotalCards = %d, want %d", cardsDrawn, s.TotalCards, initialTotal-cardsDrawn)
	}

	if len(s.Cards) != initialTotal-cardsDrawn {
		t.Errorf("Draw() %d times: Cards length = %d, want %d", cardsDrawn, len(s.Cards), initialTotal-cardsDrawn)
	}

	actualCounts := make(map[string]int)
	for _, c := range s.Cards {
		actualCounts[c.ToString()]++
	}

	for rank, count := range s.RankCounts {
		if actualCounts[rank] != count {
			t.Errorf("Draw() %d times: RankCounts[%s] = %d, but actual count is %d", cardsDrawn, rank, count, actualCounts[rank])
		}
	}
}
