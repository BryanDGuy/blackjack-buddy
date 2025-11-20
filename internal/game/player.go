package game

import (
	"errors"

	"github.com/bryan/blackjack-buddy/internal/card"
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

func (p *Player) Hit(session *Session) error {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() {
		return ErrNoActiveHand
	}

	p.ActiveHand.AddCard(session.DrawCard())

	if p.ActiveHand.IsBust() {
		p.completeAndAdvance()
		return nil
	}

	if p.ActiveHand.Value() == 21 {
		p.completeAndAdvance()
		return nil
	}

	return nil
}

func (p *Player) Stand(session *Session) error {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() {
		return ErrNoActiveHand
	}

	p.completeAndAdvance()
	return nil
}

func (p *Player) Double(session *Session) error {
	if p.ActiveHand == nil || len(p.ActiveHand.Cards) != 2 {
		return ErrInvalidMove
	}

	p.ActiveHand.AddCard(session.DrawCard())

	p.completeAndAdvance()
	return nil
}

func (p *Player) Split(session *Session) error {
	if p.ActiveHand == nil || len(p.ActiveHand.Cards) != 2 {
		return ErrInvalidMove
	}

	if !p.ActiveHand.CanSplit() {
		return ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{p.ActiveHand.Cards[0], session.DrawCard()})
	second := hand.NewHand([]card.Card{p.ActiveHand.Cards[1], session.DrawCard()})

	p.ActiveHand = first
	p.InactiveHands = append([]*hand.Hand{second}, p.InactiveHands...)

	return nil
}

func (p *Player) completeAndAdvance() {
	activeHandCardsCopy := make([]card.Card, len(p.ActiveHand.Cards))
	copy(activeHandCardsCopy, p.ActiveHand.Cards)

	p.CompletedHands = append(p.CompletedHands, hand.NewHand(activeHandCardsCopy))
	p.ActiveHand = nil

	if len(p.InactiveHands) > 0 {
		firstInactiveHandCardsCopy := make([]card.Card, len(p.InactiveHands[0].Cards))
		copy(firstInactiveHandCardsCopy, p.InactiveHands[0].Cards)
		p.ActiveHand = hand.NewHand(firstInactiveHandCardsCopy)
		p.InactiveHands = p.InactiveHands[1:]
	}
}

func (p *Player) CanMove() bool {
	return p.ActiveHand != nil && !p.ActiveHand.IsEmpty() && !p.ActiveHand.IsBust()
}
