package game

import (
	"errors"
	"math/rand"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/deck"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/player"
	"github.com/bryan/blackjack-buddy/internal/shoe"
	"github.com/bryan/blackjack-buddy/internal/strategy"
	"github.com/google/uuid"
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
	OutcomeNone      Outcome = ""
	OutcomeBust      Outcome = "Bust"
	OutcomeWin       Outcome = "Win"
	OutcomeLose      Outcome = "Lose"
	OutcomePush      Outcome = "Push"
	OutcomeBlackjack Outcome = "Blackjack"
	OutcomePending   Outcome = "Pending"
)

var (
	ErrInvalidMove  = errors.New("invalid move")
	ErrInvalidSplit = errors.New("invalid split")
	ErrNoActiveHand = errors.New("no active hand")
)

type Game struct {
	ID         string
	RoundState RoundStateType
	Player     *player.Player
	Dealer     *dealer.Dealer
	Outcomes   []Outcome
	shoe       shoe.Shoe
}

func NewGame(p *player.Player, dealer *dealer.Dealer) *Game {
	return &Game{
		ID:         uuid.New().String(),
		RoundState: RoundStateNone,
		Player:     p,
		Dealer:     dealer,
		Outcomes:   nil,
		shoe:       generateShuffledShoe(),
	}
}

func (g *Game) StartRound() {
	if g.Player == nil || g.Dealer == nil {
		return
	}

	if len(g.shoe.Cards) < reshuffleThreshold {
		g.shoe = generateShuffledShoe()
	}

	g.RoundState = RoundStateActive
	g.Player.RefreshHand(hand.NewHand([]card.Card{g.shoe.Draw(), g.shoe.Draw()}))
	g.Dealer.RefreshHand(hand.NewHand([]card.Card{g.shoe.Draw()}))
	g.Outcomes = nil
}

func (g *Game) ApplyMove(move strategy.Decision) error {
	if g.Player == nil || g.Player.ActiveHand == nil || g.Player.ActiveHand.IsEmpty() {
		return ErrNoActiveHand
	}

	var (
		isPlayerHandResolved bool  = false
		err                  error = nil
	)

	switch move {
	case strategy.Hit:
		isPlayerHandResolved = g.hit()
	case strategy.Stand:
		isPlayerHandResolved = g.stand()
	case strategy.DoubleDown:
		isPlayerHandResolved, err = g.double()
	case strategy.Split:
		isPlayerHandResolved, err = g.split()
	default:
		err = ErrInvalidMove
	}

	if err != nil {
		return err
	}

	if isPlayerHandResolved {
		g.markPlayerHandAsResolved()

		if len(g.Player.UnresolvedHands) > 0 {
			g.activateNextHand()
		} else {
			g.completeDealerHand()
			g.setOutcomes()
			g.RoundState = RoundStateComplete
		}
	}

	return nil
}

func (g *Game) hit() bool {
	g.Player.ActiveHand.AddCard(g.shoe.Draw())

	if g.Player.ActiveHand.IsBust() {
		return true
	}

	if g.Player.ActiveHand.Value() == 21 {
		return true
	}

	return false
}

func (g *Game) stand() bool {
	return true
}

func (g *Game) double() (bool, error) {
	g.Player.ActiveHand.AddCard(g.shoe.Draw())
	return true, nil
}

func (g *Game) split() (bool, error) {
	if !g.Player.ActiveHand.CanSplit() {
		return false, ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{g.Player.ActiveHand.Cards[0], g.shoe.Draw()})
	second := hand.NewHand([]card.Card{g.Player.ActiveHand.Cards[1], g.shoe.Draw()})

	g.Player.ActiveHand = first
	g.Player.UnresolvedHands = append([]*hand.Hand{second}, g.Player.UnresolvedHands...)

	return false, nil
}

func generateShuffledShoe() shoe.Shoe {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	decks := make([]deck.Deck, 0, decksInShoe)
	for range decksInShoe {
		d := deck.NewDeck(rng)
		d.Shuffle(rng)
		decks = append(decks, d)
	}

	s := shoe.NewShoe(decks)
	s.Shuffle(rng)

	return s
}

func (g *Game) markPlayerHandAsResolved() {
	activeHandCardsCopy := make([]card.Card, len(g.Player.ActiveHand.Cards))
	copy(activeHandCardsCopy, g.Player.ActiveHand.Cards)

	g.Player.ResolvedHands = append(g.Player.ResolvedHands, hand.NewHand(activeHandCardsCopy))
	g.Player.ActiveHand = nil
}

func (g *Game) activateNextHand() {
	firstUnresolvedHandCardsCopy := make([]card.Card, len(g.Player.UnresolvedHands[0].Cards))
	copy(firstUnresolvedHandCardsCopy, g.Player.UnresolvedHands[0].Cards)
	g.Player.ActiveHand = hand.NewHand(firstUnresolvedHandCardsCopy)
	g.Player.UnresolvedHands = g.Player.UnresolvedHands[1:]
}

func (g *Game) completeDealerHand() {
	if g.Dealer == nil || g.Dealer.Hand == nil {
		return
	}

	for g.Dealer.Hand.Value() < 17 {
		g.Dealer.Hand.AddCard(g.shoe.Draw())
	}
}

func (g *Game) setOutcomes() {
	if g.Dealer == nil || g.Player == nil || g.Dealer.Hand == nil {
		g.Outcomes = []Outcome{}
		return
	}

	resolvedHands := g.Player.ResolvedHands
	results := make([]Outcome, len(resolvedHands))

	for i, playerHand := range resolvedHands {
		switch {
		case playerHand.IsBust():
			results[i] = OutcomeBust
		case len(playerHand.Cards) == 2 && playerHand.Value() == 21:
			if g.Dealer.Hand.IsBlackjack() {
				results[i] = OutcomePush
			} else {
				results[i] = OutcomeBlackjack
			}
		case g.Dealer.Hand.IsBust():
			results[i] = OutcomeWin
		case playerHand.Value() > g.Dealer.Hand.Value():
			results[i] = OutcomeWin
		case playerHand.Value() < g.Dealer.Hand.Value():
			results[i] = OutcomeLose
		default:
			results[i] = OutcomePush
		}
	}

	g.Outcomes = results
}
