package game

import (
	"github.com/bryan/blackjack-buddy/internal/card"
)

type RoundStateType string

const (
	RoundStateNone     RoundStateType = "none"
	RoundStateActive   RoundStateType = "active"
	RoundStateComplete RoundStateType = "complete"
)

type GameSession struct {
	Session        *Session
	RoundState     RoundStateType
	ActiveHand     []card.Card
	InactiveHands  [][]card.Card
	CompletedHands [][]card.Card
	DealerCard     card.Card
	DealerCards    []card.Card
	Outcomes       []string
}

func NewGameSession(session *Session) *GameSession {
	return &GameSession{
		Session:        session,
		RoundState:     RoundStateNone,
		ActiveHand:     nil,
		InactiveHands:  nil,
		CompletedHands: nil,
		DealerCard:     card.Card{},
		DealerCards:    nil,
		Outcomes:       nil,
	}
}

