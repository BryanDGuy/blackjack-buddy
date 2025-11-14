package card

import "testing"

func TestCard_Value(t *testing.T) {
	tests := []struct {
		rank  Rank
		value int
	}{
		{Two, 2},
		{Three, 3},
		{Four, 4},
		{Five, 5},
		{Six, 6},
		{Seven, 7},
		{Eight, 8},
		{Nine, 9},
		{Ten, 10},
		{Jack, 10},
		{Queen, 10},
		{King, 10},
		{Ace, 11},
	}

	for _, tt := range tests {
		c := NewCard(tt.rank)
		if got := c.Value(); got != tt.value {
			t.Errorf("Card(%v).Value() = %d, want %d", tt.rank, got, tt.value)
		}
	}
}

func TestCard_ToString(t *testing.T) {
	tests := []struct {
		rank Rank
		want string
	}{
		{Two, "2"},
		{Three, "3"},
		{Four, "4"},
		{Five, "5"},
		{Six, "6"},
		{Seven, "7"},
		{Eight, "8"},
		{Nine, "9"},
		{Ten, "10"},
		{Jack, "J"},
		{Queen, "Q"},
		{King, "K"},
		{Ace, "A"},
	}

	for _, tt := range tests {
		c := NewCard(tt.rank)
		if got := c.ToString(); got != tt.want {
			t.Errorf("Card(%v).ToString() = %q, want %q", tt.rank, got, tt.want)
		}
	}
}
