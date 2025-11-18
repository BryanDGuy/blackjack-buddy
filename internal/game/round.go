package game

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type RoundState struct {
	Player    []card.Card
	Dealer    card.Card
	Queue     [][]card.Card
	Completed [][]card.Card
	Outcomes  []Outcome
}

type RoundResolution struct {
	State         RoundState
	Outcome       string
	RoundComplete bool
	DealerCards   []card.Card
}

func NewRoundState(player []card.Card, dealer card.Card, queued [][]card.Card, completed [][]card.Card, outcomes []Outcome) RoundState {
	return RoundState{
		Player:    cloneCards(player),
		Dealer:    dealer,
		Queue:     cloneHands(queued),
		Completed: cloneHands(completed),
		Outcomes:  append([]Outcome{}, outcomes...),
	}
}

func ApplyDecision(state RoundState, decision strategy.Decision, session *Session) (RoundResolution, error) {
	switch decision {
	case strategy.Hit:
		return applyHit(state, session), nil
	case strategy.DoubleDown:
		return applyDouble(state, session), nil
	case strategy.Stand:
		return applyStand(state, session), nil
	case strategy.Split:
		return applySplit(state, session)
	default:
		return RoundResolution{}, errors.New("unsupported decision")
	}
}

func applyHit(state RoundState, session *Session) RoundResolution {
	res := newResolution(state)
	res.State.Player = append(res.State.Player, session.DrawCard())

	playerHand := hand.NewHand(res.State.Player)
	if playerHand.IsBust() {
		res.Outcome = OutcomeBust.String()
		res.RoundComplete = true
	} else if playerHand.Value() == 21 {
		res.RoundComplete = true
	}

	return finalize(res, session)
}

func applyDouble(state RoundState, session *Session) RoundResolution {
	res := newResolution(state)
	res.State.Player = append(res.State.Player, session.DrawCard())
	res.RoundComplete = true

	if hand.NewHand(res.State.Player).IsBust() {
		res.Outcome = OutcomeBust.String()
	}

	return finalize(res, session)
}

func applyStand(state RoundState, session *Session) RoundResolution {
	res := newResolution(state)
	res.RoundComplete = true
	return finalize(res, session)
}

func applySplit(state RoundState, session *Session) (RoundResolution, error) {
	if !hand.NewHand(state.Player).CanSplit() {
		return RoundResolution{}, ErrInvalidSplit
	}

	res := newResolution(state)
	first := []card.Card{res.State.Player[0], session.DrawCard()}
	second := []card.Card{res.State.Player[1], session.DrawCard()}

	res.State.Player = first
	res.State.Queue = append([][]card.Card{second}, res.State.Queue...)
	res.DealerCards = []card.Card{res.State.Dealer}

	return res, nil
}

func finalize(res RoundResolution, session *Session) RoundResolution {
	if !res.RoundComplete {
		return res
	}

	outcome := parseOutcomeFromString(res.Outcome)
	res.State = appendCompleted(res.State, outcome)

	if len(res.State.Queue) > 0 {
		res = advanceQueue(res)
	} else {
		res = settleDealer(res, session)
	}

	return res
}

func newResolution(state RoundState) RoundResolution {
	return RoundResolution{
		State: RoundState{
			Player:    cloneCards(state.Player),
			Dealer:    state.Dealer,
			Queue:     cloneHands(state.Queue),
			Completed: cloneHands(state.Completed),
			Outcomes:  append([]Outcome{}, state.Outcomes...),
		},
		DealerCards: []card.Card{state.Dealer},
	}
}

func appendCompleted(state RoundState, outcome Outcome) RoundState {
	current := cloneCards(state.Player)
	state.Completed = append(state.Completed, current)

	if outcome == OutcomeNone && hand.NewHand(current).IsBust() {
		outcome = OutcomeBust
	}

	state.Outcomes = append(state.Outcomes, outcome)
	return state
}

func advanceQueue(res RoundResolution) RoundResolution {
	next := cloneCards(res.State.Queue[0])
	res.State.Queue = cloneHands(res.State.Queue[1:])
	res.State.Player = next
	res.RoundComplete = false
	res.Outcome = ""
	res.DealerCards = []card.Card{res.State.Dealer}
	return res
}

func settleDealer(res RoundResolution, session *Session) RoundResolution {
	if session.Dealer() != nil && session.Dealer().Hand() != nil && session.Dealer().Hand().Count() > 0 {
		for session.Dealer().Hand().Value() < 17 {
			session.Dealer().Hand().AddCard(session.DrawCard())
		}
	}

	outcomes := session.DetermineOutcome()

	var dealerCards []card.Card
	if session.Dealer() != nil && session.Dealer().Hand() != nil {
		dealerCards = session.Dealer().Hand().Cards()
	}
	res.DealerCards = dealerCards

	for i := range res.State.Outcomes {
		if res.State.Outcomes[i] == OutcomeNone && i < len(outcomes) {
			res.State.Outcomes[i] = outcomes[i]
		}
	}

	for i := len(res.State.Outcomes); i < len(outcomes); i++ {
		res.State.Outcomes = append(res.State.Outcomes, outcomes[i])
	}

	if len(outcomes) == 0 {
		res.Outcome = ""
	} else {
		parts := make([]string, 0, len(outcomes))
		for i, outcome := range outcomes {
			label := outcome.String()
			if label == "" {
				label = OutcomePending.String()
			}
			parts = append(parts, fmt.Sprintf("Hand%d %s", i+1, label))
		}
		res.Outcome = strings.Join(parts, " | ")
	}

	return res
}

func cloneCards(cards []card.Card) []card.Card {
	out := make([]card.Card, len(cards))
	copy(out, cards)
	return out
}

func cloneHands(hands [][]card.Card) [][]card.Card {
	out := make([][]card.Card, len(hands))
	for i, h := range hands {
		out[i] = cloneCards(h)
	}
	return out
}

func parseOutcomeFromString(s string) Outcome {
	if s == "" {
		return OutcomeNone
	}
	return ParseOutcome(s)
}
