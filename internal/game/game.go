package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/deck"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/player"
	"github.com/bryan/blackjack-buddy/internal/shoe"
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

type Game struct {
	ID         string
	Shoe       shoe.Shoe
	RoundState RoundStateType
	Player     *player.Player
	Dealer     *dealer.Dealer
	Outcomes   []Outcome
}

func NewGame(p *player.Player, dealer *dealer.Dealer) *Game {
	return &Game{
		ID:         uuid.New().String(),
		Shoe:       generateShuffledShoe(),
		RoundState: RoundStateNone,
		Player:     p,
		Dealer:     dealer,
		Outcomes:   nil,
	}
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
	(&s).Shuffle(rng)

	return s
}

func (g *Game) DrawCard() card.Card {
	if len(g.Shoe.Cards) < reshuffleThreshold {
		g.Shoe = generateShuffledShoe()
	}
	return g.Shoe.Draw()
}

func (g *Game) Hit() error {
	if g.Player == nil || g.Player.ActiveHand == nil || g.Player.ActiveHand.IsEmpty() {
		return player.ErrNoActiveHand
	}

	g.Player.ActiveHand.AddCard(g.DrawCard())

	if g.Player.ActiveHand.IsBust() {
		g.completeAndAdvance()
		return nil
	}

	if g.Player.ActiveHand.Value() == 21 {
		g.completeAndAdvance()
		return nil
	}

	return nil
}

func (g *Game) Stand() error {
	if g.Player == nil || g.Player.ActiveHand == nil || g.Player.ActiveHand.IsEmpty() {
		return player.ErrNoActiveHand
	}

	g.completeAndAdvance()
	return nil
}

func (g *Game) Double() error {
	if g.Player == nil || g.Player.ActiveHand == nil || len(g.Player.ActiveHand.Cards) != 2 {
		return player.ErrInvalidMove
	}

	g.Player.ActiveHand.AddCard(g.DrawCard())
	g.completeAndAdvance()
	return nil
}

func (g *Game) Split() error {
	if g.Player == nil || g.Player.ActiveHand == nil || len(g.Player.ActiveHand.Cards) != 2 {
		return player.ErrInvalidMove
	}

	if !g.Player.ActiveHand.CanSplit() {
		return player.ErrInvalidSplit
	}

	first := hand.NewHand([]card.Card{g.Player.ActiveHand.Cards[0], g.DrawCard()})
	second := hand.NewHand([]card.Card{g.Player.ActiveHand.Cards[1], g.DrawCard()})

	g.Player.ActiveHand = first
	g.Player.InactiveHands = append([]*hand.Hand{second}, g.Player.InactiveHands...)

	return nil
}

func (g *Game) completeAndAdvance() {
	if g.Player == nil || g.Player.ActiveHand == nil {
		return
	}

	activeHandCardsCopy := make([]card.Card, len(g.Player.ActiveHand.Cards))
	copy(activeHandCardsCopy, g.Player.ActiveHand.Cards)

	g.Player.CompletedHands = append(g.Player.CompletedHands, hand.NewHand(activeHandCardsCopy))
	g.Player.ActiveHand = nil

	if len(g.Player.InactiveHands) > 0 {
		firstInactiveHandCardsCopy := make([]card.Card, len(g.Player.InactiveHands[0].Cards))
		copy(firstInactiveHandCardsCopy, g.Player.InactiveHands[0].Cards)
		g.Player.ActiveHand = hand.NewHand(firstInactiveHandCardsCopy)
		g.Player.InactiveHands = g.Player.InactiveHands[1:]
	}
}

func (g *Game) DetermineOutcome() []Outcome {
	if g.Dealer == nil || g.Player == nil || g.Dealer.Hand == nil {
		return nil
	}

	completedHands := g.Player.CompletedHands
	results := make([]Outcome, len(completedHands))

	for i, playerHand := range completedHands {
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

	return results
}

func (g *Game) UpdateOutcomes(newOutcomes []Outcome) {
	for i := range g.Outcomes {
		if g.Outcomes[i] == OutcomeNone && i < len(newOutcomes) {
			g.Outcomes[i] = newOutcomes[i]
		}
	}

	for i := len(g.Outcomes); i < len(newOutcomes); i++ {
		g.Outcomes = append(g.Outcomes, newOutcomes[i])
	}
}

func (g *Game) FormatOutcomes() string {
	if len(g.Outcomes) == 0 {
		return ""
	}

	parts := make([]string, 0, len(g.Outcomes))
	for i, outcome := range g.Outcomes {
		label := string(outcome)
		if label == "" {
			label = string(OutcomePending)
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}

	return strings.Join(parts, " | ")
}
