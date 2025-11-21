package player

import (
	"errors"

	"github.com/bryan/blackjack-buddy/internal/hand"
)

var ErrInvalidMove = errors.New("invalid move")
var ErrInvalidSplit = errors.New("invalid split")
var ErrNoActiveHand = errors.New("no active hand")

type Player struct {
	ActiveHand      *hand.Hand
	UnresolvedHands []*hand.Hand
	ResolvedHands   []*hand.Hand
}

func NewPlayer() *Player {
	return &Player{
		ActiveHand:      nil,
		UnresolvedHands: nil,
		ResolvedHands:   nil,
	}
}

func (p *Player) CanMove() bool {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() || p.ActiveHand.IsBust() {
		return false
	}
	if p.ActiveHand.Value() == 21 {
		return false
	}
	return true
}
