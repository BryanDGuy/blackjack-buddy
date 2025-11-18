package game

import (
	"errors"
	"fmt"
	"strings"

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
	Outcomes       []string
}

type PlayerMoveResult struct {
	RoundComplete bool
	Outcome       string
}

func NewPlayer(hand *hand.Hand) *Player {
	return &Player{
		ActiveHand:     hand,
		InactiveHands:  nil,
		CompletedHands: nil,
		Outcomes:       nil,
	}
}

func (p *Player) Hit(session *Session) (PlayerMoveResult, error) {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() {
		return PlayerMoveResult{}, ErrNoActiveHand
	}

	p.ActiveHand.Cards = append(p.ActiveHand.Cards, session.DrawCard())

	result := PlayerMoveResult{}
	if p.ActiveHand.IsBust() {
		result.Outcome = "Bust"
		result.RoundComplete = true
		p.completeActiveHand(result.Outcome)
		return p.finalizeOrAdvance(result)
	}

	if p.ActiveHand.Value() == 21 {
		result.RoundComplete = true
		p.completeActiveHand("")
		return p.finalizeOrAdvance(result)
	}

	return result, nil
}

func (p *Player) Stand(session *Session) (PlayerMoveResult, error) {
	if p.ActiveHand == nil || p.ActiveHand.IsEmpty() {
		return PlayerMoveResult{}, ErrNoActiveHand
	}

	result := PlayerMoveResult{RoundComplete: true}
	p.completeActiveHand("")
	return p.finalizeOrAdvance(result)
}

func (p *Player) Double(session *Session) (PlayerMoveResult, error) {
	if p.ActiveHand == nil || len(p.ActiveHand.Cards) != 2 {
		return PlayerMoveResult{}, ErrInvalidMove
	}

	p.ActiveHand.Cards = append(p.ActiveHand.Cards, session.DrawCard())
	result := PlayerMoveResult{RoundComplete: true}

	if p.ActiveHand.IsBust() {
		result.Outcome = "Bust"
	}

	p.completeActiveHand(result.Outcome)
	return p.finalizeOrAdvance(result)
}

func (p *Player) Split(session *Session) (PlayerMoveResult, error) {
	if p.ActiveHand == nil || len(p.ActiveHand.Cards) != 2 {
		return PlayerMoveResult{}, ErrInvalidMove
	}

	if !p.ActiveHand.CanSplit() {
		return PlayerMoveResult{}, ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{p.ActiveHand.Cards[0], session.DrawCard()})
	second := hand.NewHand([]card.Card{p.ActiveHand.Cards[1], session.DrawCard()})

	p.ActiveHand = first
	p.InactiveHands = append([]*hand.Hand{second}, p.InactiveHands...)

	return PlayerMoveResult{RoundComplete: false}, nil
}

func (p *Player) completeActiveHand(outcome string) {
	completed := hand.NewHand(p.ActiveHand.Cards)
	p.CompletedHands = append(p.CompletedHands, completed)

	if outcome == "" {
		if completed.IsBust() {
			outcome = "Bust"
		}
	}

	p.Outcomes = append(p.Outcomes, outcome)
}

func (p *Player) finalizeOrAdvance(result PlayerMoveResult) (PlayerMoveResult, error) {
	if len(p.InactiveHands) > 0 {
		p.ActiveHand = hand.NewHand(p.InactiveHands[0].Cards)
		p.InactiveHands = p.InactiveHands[1:]
		result.RoundComplete = false
		result.Outcome = ""
		return result, nil
	}

	return result, nil
}

func (p *Player) Finalize(dealer *Dealer) {
	if len(p.CompletedHands) == 0 || dealer == nil {
		return
	}

	outcomes := dealer.EvaluatePlayerHands(p.CompletedHands)
	p.updateOutcomes(outcomes)
}

func (p *Player) updateOutcomes(newOutcomes []string) {
	for i := range p.Outcomes {
		if p.Outcomes[i] == "" && i < len(newOutcomes) {
			p.Outcomes[i] = newOutcomes[i]
		}
	}

	for i := len(p.Outcomes); i < len(newOutcomes); i++ {
		p.Outcomes = append(p.Outcomes, newOutcomes[i])
	}
}

func (p *Player) FormatOutcomes() string {
	if len(p.Outcomes) == 0 {
		return ""
	}

	parts := make([]string, 0, len(p.Outcomes))
	for i, outcome := range p.Outcomes {
		label := outcome
		if label == "" {
			label = "Pending"
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}

	return strings.Join(parts, " | ")
}
