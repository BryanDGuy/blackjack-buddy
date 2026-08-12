package player

import "testing"

func TestNewPlayer(t *testing.T) {
	p := NewPlayer()
	if p == nil {
		t.Fatal("NewPlayer() returned nil")
	}
	if p.ActiveHand != nil {
		t.Error("NewPlayer() ActiveHand should be nil")
	}
	if p.UnresolvedHands != nil {
		t.Error("NewPlayer() UnresolvedHands should be nil")
	}
	if p.ResolvedHands != nil {
		t.Error("NewPlayer() ResolvedHands should be nil")
	}
}
