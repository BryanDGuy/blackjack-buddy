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
	activeHand     *hand.Hand
	inactiveHands  []*hand.Hand
	completedHands []*hand.Hand
}

func NewPlayer() *Player {
	return &Player{
		activeHand:     nil,
		inactiveHands:  nil,
		completedHands: nil,
	}
}

func (p *Player) ActiveHand() *hand.Hand {
	return p.activeHand
}

func (p *Player) SetActiveHand(h *hand.Hand) {
	p.activeHand = h
}

func (p *Player) InactiveHands() []*hand.Hand {
	hands := make([]*hand.Hand, len(p.inactiveHands))
	copy(hands, p.inactiveHands)
	return hands
}

func (p *Player) SetInactiveHands(hands []*hand.Hand) {
	p.inactiveHands = make([]*hand.Hand, len(hands))
	copy(p.inactiveHands, hands)
}

func (p *Player) CompletedHands() []*hand.Hand {
	hands := make([]*hand.Hand, len(p.completedHands))
	copy(hands, p.completedHands)
	return hands
}

func (p *Player) SetCompletedHands(hands []*hand.Hand) {
	p.completedHands = make([]*hand.Hand, len(hands))
	copy(p.completedHands, hands)
}

func (p *Player) Hit(session *Session) error {
	if p.activeHand == nil || p.activeHand.IsEmpty() {
		return ErrNoActiveHand
	}

	p.activeHand.AddCard(session.DrawCard())

	if p.activeHand.IsBust() {
		p.completeAndAdvance()
		return nil
	}

	if p.activeHand.Value() == 21 {
		p.completeAndAdvance()
		return nil
	}

	return nil
}

func (p *Player) Stand(session *Session) error {
	if p.activeHand == nil || p.activeHand.IsEmpty() {
		return ErrNoActiveHand
	}

	p.completeAndAdvance()
	return nil
}

func (p *Player) Double(session *Session) error {
	if p.activeHand == nil || p.activeHand.Count() != 2 {
		return ErrInvalidMove
	}

	p.activeHand.AddCard(session.DrawCard())

	p.completeAndAdvance()
	return nil
}

func (p *Player) Split(session *Session) error {
	if p.activeHand == nil || p.activeHand.Count() != 2 {
		return ErrInvalidMove
	}

	if !p.activeHand.CanSplit() {
		return ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{p.activeHand.GetCard(0), session.DrawCard()})
	second := hand.NewHand([]card.Card{p.activeHand.GetCard(1), session.DrawCard()})

	p.activeHand = first
	p.inactiveHands = append([]*hand.Hand{second}, p.inactiveHands...)

	return nil
}

func (p *Player) completeAndAdvance() {
	completed := hand.NewHand(p.activeHand.Cards())
	p.completedHands = append(p.completedHands, completed)

	if len(p.inactiveHands) > 0 {
		p.activeHand = hand.NewHand(p.inactiveHands[0].Cards())
		p.inactiveHands = p.inactiveHands[1:]
	}
}

func (p *Player) CanMove() bool {
	return p.activeHand != nil && !p.activeHand.IsEmpty() && !p.activeHand.IsBust()
}
