package player

import (
	"errors"

	"github.com/bryan/blackjack-buddy/internal/hand"
)

var ErrInvalidMove = errors.New("invalid move")
var ErrInvalidSplit = errors.New("invalid split")
var ErrNoActiveHand = errors.New("no active hand")

type Player struct {
	ActiveHand     *hand.Hand
	InactiveHands  []*hand.Hand
	CompletedHands []*hand.Hand
}

func NewPlayer() *Player {
	return &Player{
		ActiveHand:     nil,
		InactiveHands:  nil,
		CompletedHands: nil,
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
