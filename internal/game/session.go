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
	id         string
	shoe       shoe.Shoe
	roundState RoundStateType
	player     *Player
	dealer     *dealer.Dealer
	outcomes   []Outcome
}

func NewSession(player *Player, dealer *dealer.Dealer) *Session {
	return &Session{
		id:         uuid.New().String(),
		shoe:       generateShuffledShoe(),
		roundState: RoundStateNone,
		player:     player,
		dealer:     dealer,
		outcomes:   nil,
	}
}

func (s *Session) Player() *Player {
	return s.player
}

func (s *Session) SetPlayer(p *Player) {
	s.player = p
}

func (s *Session) Dealer() *dealer.Dealer {
	return s.dealer
}

func (s *Session) SetDealer(d *dealer.Dealer) {
	s.dealer = d
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) Shoe() shoe.Shoe {
	return s.shoe
}

func (s *Session) SetShoe(sh shoe.Shoe) {
	s.shoe = sh
}

func (s *Session) RoundState() RoundStateType {
	return s.roundState
}

func (s *Session) SetRoundState(state RoundStateType) {
	s.roundState = state
}

func (s *Session) Outcomes() []Outcome {
	outcomes := make([]Outcome, len(s.outcomes))
	copy(outcomes, s.outcomes)
	return outcomes
}

func (s *Session) SetOutcomes(outcomes []Outcome) {
	s.outcomes = make([]Outcome, len(outcomes))
	copy(s.outcomes, outcomes)
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

func (s *Session) DrawCard() card.Card {
	if len(s.shoe.Cards()) < reshuffleThreshold {
		s.shoe = generateShuffledShoe()
	}
	return s.shoe.Draw()
}

func (s *Session) DetermineOutcome() []Outcome {
	if s.dealer == nil || s.player == nil || s.dealer.Hand() == nil {
		return nil
	}

	completedHands := s.player.CompletedHands()
	results := make([]Outcome, len(completedHands))

	for i, playerHand := range completedHands {
		switch {
		case playerHand.IsBust():
			results[i] = OutcomeBust
		case playerHand.Count() == 2 && playerHand.Value() == 21:
			if s.dealer.Hand().IsBlackjack() {
				results[i] = OutcomePush
			} else {
				results[i] = OutcomeBlackjack
			}
		case s.dealer.Hand().IsBust():
			results[i] = OutcomeWin
		case playerHand.Value() > s.dealer.Hand().Value():
			results[i] = OutcomeWin
		case playerHand.Value() < s.dealer.Hand().Value():
			results[i] = OutcomeLose
		default:
			results[i] = OutcomePush
		}
	}

	return results
}

func (s *Session) UpdateOutcomes(newOutcomes []Outcome) {
	for i := range s.outcomes {
		if s.outcomes[i] == OutcomeNone && i < len(newOutcomes) {
			s.outcomes[i] = newOutcomes[i]
		}
	}

	for i := len(s.outcomes); i < len(newOutcomes); i++ {
		s.outcomes = append(s.outcomes, newOutcomes[i])
	}
}

func (s *Session) FormatOutcomes() string {
	if len(s.outcomes) == 0 {
		return ""
	}

	parts := make([]string, 0, len(s.outcomes))
	for i, outcome := range s.outcomes {
		label := string(outcome)
		if label == "" {
			label = string(OutcomePending)
		}
		parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
	}

	return strings.Join(parts, " | ")
}
