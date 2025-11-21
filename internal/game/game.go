package game

import (
	"math/rand"
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
	g.Player.UnresolvedHands = append([]*hand.Hand{second}, g.Player.UnresolvedHands...)

	return nil
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

func (g *Game) completeAndAdvance() {
	if g.Player == nil || g.Player.ActiveHand == nil {
		return
	}

	g.completePlayerHand()

	if len(g.Player.UnresolvedHands) > 0 {
		g.activateNextHand()
	} else {
		g.setOutcomes()
	}
}

func (g *Game) completePlayerHand() {
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
