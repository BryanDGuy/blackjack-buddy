// Package player manages player hands.
package player

import "github.com/bryan/blackjack-buddy/internal/hand"

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

func (p *Player) RefreshHand(newHand *hand.Hand) {
	p.ActiveHand = newHand
	p.UnresolvedHands = nil
	p.ResolvedHands = nil
}
