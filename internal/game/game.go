// Package game manages blackjack rounds.
package game

import (
	cryptorand "crypto/rand"
	"errors"
	"math/rand"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/shoe"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

const (
	decksInShoe        = 6
	reshuffleThreshold = 78
)

type RoundStateType string

const (
	RoundStateNone     RoundStateType = "none"
	RoundStateActive   RoundStateType = "active"
	RoundStateComplete RoundStateType = "complete"
)

type Outcome string

const (
	OutcomeBust      Outcome = "Bust"
	OutcomeWin       Outcome = "Win"
	OutcomeLose      Outcome = "Lose"
	OutcomePush      Outcome = "Push"
	OutcomeBlackjack Outcome = "Blackjack"
)

var (
	ErrInvalidMove  = errors.New("invalid move")
	ErrInvalidSplit = errors.New("invalid split")
	ErrNoActiveHand = errors.New("no active hand")
)

type Game struct {
	ID              string
	RoundState      RoundStateType
	ActiveHand      *hand.Hand
	UnresolvedHands []*hand.Hand
	ResolvedHands   []*hand.Hand
	DealerHand      *hand.Hand
	Outcomes        []Outcome
	shoe            shoe.Shoe
}

func NewGame() *Game {
	return &Game{
		ID:         cryptorand.Text(),
		RoundState: RoundStateNone,
		shoe:       generateShuffledShoe(),
	}
}

func (g *Game) StartRound() {
	if len(g.shoe.Cards) < reshuffleThreshold {
		g.shoe = generateShuffledShoe()
	}

	g.RoundState = RoundStateActive
	g.ActiveHand = hand.NewHand([]card.Card{g.shoe.Draw(), g.shoe.Draw()})
	g.UnresolvedHands = nil
	g.ResolvedHands = nil
	g.DealerHand = hand.NewHand([]card.Card{g.shoe.Draw(), g.shoe.Draw()})
	g.Outcomes = nil
}

func (g *Game) AbandonRound() {
	if g.RoundState != RoundStateActive {
		return
	}
	g.ActiveHand = nil
	g.UnresolvedHands = nil
	g.ResolvedHands = nil
	g.Outcomes = nil
	g.RoundState = RoundStateComplete
}

func (g *Game) ApplyMove(move strategy.Decision) error {
	if g.ActiveHand == nil || g.ActiveHand.IsEmpty() {
		return ErrNoActiveHand
	}

	var (
		isPlayerHandResolved = false
		err                  error
	)

	switch move {
	case strategy.Hit:
		isPlayerHandResolved = g.hit()
	case strategy.Stand:
		isPlayerHandResolved = true
	case strategy.DoubleDown:
		isPlayerHandResolved, err = g.double()
	case strategy.Split:
		err = g.split()
	default:
		err = ErrInvalidMove
	}

	if err != nil {
		return err
	}

	if isPlayerHandResolved {
		g.ResolvedHands = append(g.ResolvedHands, g.ActiveHand)
		g.ActiveHand = nil

		if len(g.UnresolvedHands) > 0 {
			g.ActiveHand = g.UnresolvedHands[0]
			g.UnresolvedHands = g.UnresolvedHands[1:]
		} else {
			g.completeDealerHand()
			g.setOutcomes()
			g.RoundState = RoundStateComplete
		}
	}

	return nil
}

func (g *Game) hit() bool {
	g.ActiveHand.AddCard(g.shoe.Draw())
	return g.ActiveHand.Value() >= 21
}

func (g *Game) double() (bool, error) {
	if len(g.ActiveHand.Cards) != 2 {
		return false, ErrInvalidMove
	}
	g.ActiveHand.AddCard(g.shoe.Draw())
	return true, nil
}

func (g *Game) split() error {
	if !g.ActiveHand.CanSplit() {
		return ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{g.ActiveHand.Cards[0], g.shoe.Draw()})
	second := hand.NewHand([]card.Card{g.ActiveHand.Cards[1], g.shoe.Draw()})
	first.FromSplit = true
	second.FromSplit = true

	g.ActiveHand = first
	g.UnresolvedHands = append([]*hand.Hand{second}, g.UnresolvedHands...)

	return nil
}

func generateShuffledShoe() shoe.Shoe {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := shoe.NewShoe(decksInShoe)
	s.Shuffle(rng)

	return s
}

func (g *Game) completeDealerHand() {
	if g.DealerHand == nil {
		return
	}

	for g.DealerHand.Value() < 17 {
		g.DealerHand.AddCard(g.shoe.Draw())
	}
}

func (g *Game) setOutcomes() {
	if g.DealerHand == nil {
		g.Outcomes = []Outcome{}
		return
	}

	resolvedHands := g.ResolvedHands
	results := make([]Outcome, len(resolvedHands))

	for i, playerHand := range resolvedHands {
		switch {
		case playerHand.IsBust():
			results[i] = OutcomeBust
		case !playerHand.FromSplit && len(playerHand.Cards) == 2 && playerHand.Value() == 21:
			if g.DealerHand.IsBlackjack() {
				results[i] = OutcomePush
			} else {
				results[i] = OutcomeBlackjack
			}
		case g.DealerHand.IsBlackjack():
			results[i] = OutcomeLose
		case g.DealerHand.IsBust():
			results[i] = OutcomeWin
		case playerHand.Value() > g.DealerHand.Value():
			results[i] = OutcomeWin
		case playerHand.Value() < g.DealerHand.Value():
			results[i] = OutcomeLose
		default:
			results[i] = OutcomePush
		}
	}

	g.Outcomes = results
}
