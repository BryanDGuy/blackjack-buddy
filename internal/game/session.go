package game

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/deck"
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

type Session struct {
	ID         string
	Shoe       shoe.Shoe
	RoundState RoundStateType
	Player     *Player
	Dealer     *dealer.Dealer
	Outcomes   []Outcome
}

func NewSession(player *Player, dealer *dealer.Dealer) *Session {
	return &Session{
		ID:         uuid.New().String(),
		Shoe:       generateShuffledShoe(),
		RoundState: RoundStateNone,
		Player:     player,
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

func (s *Session) DrawCard() card.Card {
	if len(s.Shoe.Cards) < reshuffleThreshold {
		s.Shoe = generateShuffledShoe()
	}
	return s.Shoe.Draw()
}

func (s *Session) DetermineOutcome() []Outcome {
	if s.Dealer == nil || s.Player == nil || s.Dealer.Hand == nil {
		return nil
	}

	completedHands := s.Player.CompletedHands
	results := make([]Outcome, len(completedHands))

	for i, playerHand := range completedHands {
		switch {
		case playerHand.IsBust():
			results[i] = OutcomeBust
		case len(playerHand.Cards) == 2 && playerHand.Value() == 21:
			if s.Dealer.Hand.IsBlackjack() {
				results[i] = OutcomePush
			} else {
				results[i] = OutcomeBlackjack
			}
		case s.Dealer.Hand.IsBust():
			results[i] = OutcomeWin
		case playerHand.Value() > s.Dealer.Hand.Value():
			results[i] = OutcomeWin
		case playerHand.Value() < s.Dealer.Hand.Value():
			results[i] = OutcomeLose
		default:
			results[i] = OutcomePush
		}
	}

	return results
}

func (s *Session) UpdateOutcomes(newOutcomes []Outcome) {
	for i := range s.Outcomes {
		if s.Outcomes[i] == OutcomeNone && i < len(newOutcomes) {
			s.Outcomes[i] = newOutcomes[i]
		}
	}

	for i := len(s.Outcomes); i < len(newOutcomes); i++ {
		s.Outcomes = append(s.Outcomes, newOutcomes[i])
	}
}

func (s *Session) FormatOutcomes() string {
	if len(s.Outcomes) == 0 {
		return ""
	}

	parts := make([]string, 0, len(s.Outcomes))
	for i, outcome := range s.Outcomes {
		label := string(outcome)
		if label == "" {
			label = string(OutcomePending)
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}

	return strings.Join(parts, " | ")
}
