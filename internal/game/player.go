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

	p.ActiveHand.Cards = append(p.ActiveHand.Cards, session.DrawCard())

	if p.ActiveHand.IsBust() {
		p.completeActiveHand(session, OutcomeBust)
		p.finalizeOrAdvance()
		return nil
	}

	if p.ActiveHand.Value() == 21 {
		p.completeActiveHand(session, OutcomeNone)
		p.finalizeOrAdvance()
		return nil
	}

	return nil
}

func (p *Player) Stand(session *Session) error {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() {
		return ErrNoActiveHand
	}

	p.completeActiveHand(session, OutcomeNone)
	p.finalizeOrAdvance()
	return nil
}

func (p *Player) Double(session *Session) error {
	if p.ActiveHand == nil || len(p.ActiveHand.Cards) != 2 {
		return ErrInvalidMove
	}

	p.ActiveHand.Cards = append(p.ActiveHand.Cards, session.DrawCard())

	outcome := OutcomeNone
	if p.ActiveHand.IsBust() {
		outcome = OutcomeBust
	}
	p.completeActiveHand(session, outcome)
	p.finalizeOrAdvance()
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

func (p *Player) completeActiveHand(session *Session, outcome Outcome) {
	completed := hand.NewHand(p.ActiveHand.Cards)
	p.CompletedHands = append(p.CompletedHands, completed)

	if outcome == OutcomeNone {
		if completed.IsBust() {
			outcome = OutcomeBust
		}
	}

	session.Outcomes = append(session.Outcomes, outcome)
}

func (p *Player) finalizeOrAdvance() {
	if len(p.InactiveHands) > 0 {
		p.ActiveHand = hand.NewHand(p.InactiveHands[0].Cards)
		p.InactiveHands = p.InactiveHands[1:]
	}
}

func (p *Player) CanMove() bool {
	return p.ActiveHand != nil && !p.ActiveHand.IsEmpty() && !p.ActiveHand.IsBust()
}
