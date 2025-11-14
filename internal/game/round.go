package game

import (
	"errors"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

var ErrInvalidSplit = errors.New("invalid split")

type RoundState struct {
	Player    []card.Card
	Dealer    card.Card
	Queue     [][]card.Card
	Completed [][]card.Card
	Outcomes  []string
	HandBets  []int
}

type RoundResolution struct {
	State          RoundState
	Outcome        string
	RoundComplete  bool
	DealerCards    []card.Card
	CurrentHandBet int
}

func NewRoundState(player []card.Card, dealer card.Card, queued [][]card.Card, completed [][]card.Card, outcomes []string) RoundState {
	return RoundState{
		Player:    cloneCards(player),
		Dealer:    dealer,
		Queue:     cloneHands(queued),
		Completed: cloneHands(completed),
		Outcomes:  append([]string{}, outcomes...),
		HandBets:  make([]int, len(completed)),
	}
}

func ApplyDecision(state RoundState, decision strategy.Decision, engine *Engine, bet int) (RoundResolution, error) {
	switch decision {
	case strategy.Hit:
		return applyHit(state, engine, bet), nil
	case strategy.DoubleDown:
		return applyDouble(state, engine, bet), nil
	case strategy.Stand:
		return applyStand(state, engine, bet), nil
	case strategy.Split:
		return applySplit(state, engine, bet)
	default:
		return RoundResolution{}, errors.New("unsupported decision")
	}
}

func applyHit(state RoundState, engine *Engine, bet int) RoundResolution {
	res := newResolution(state)
	res.State.Player = append(res.State.Player, engine.DrawCard())
	res.CurrentHandBet = bet

	playerHand := hand.NewHand(res.State.Player)
	if playerHand.IsBust() {
		res.Outcome = "Bust"
		res.RoundComplete = true
	} else if playerHand.Value() == 21 {
		res.RoundComplete = true
	}

	return finalize(res, engine)
}

func applyDouble(state RoundState, engine *Engine, bet int) RoundResolution {
	res := newResolution(state)
	res.State.Player = append(res.State.Player, engine.DrawCard())
	res.RoundComplete = true
	res.CurrentHandBet = bet * 2

	if hand.NewHand(res.State.Player).IsBust() {
		res.Outcome = "Bust"
	}

	return finalize(res, engine)
}

func applyStand(state RoundState, engine *Engine, bet int) RoundResolution {
	res := newResolution(state)
	res.RoundComplete = true
	res.CurrentHandBet = bet
	return finalize(res, engine)
}

func applySplit(state RoundState, engine *Engine, bet int) (RoundResolution, error) {
	if !hand.NewHand(state.Player).CanSplit() {
		return RoundResolution{}, ErrInvalidSplit
	}

	res := newResolution(state)
	first := []card.Card{res.State.Player[0], engine.DrawCard()}
	second := []card.Card{res.State.Player[1], engine.DrawCard()}

	res.State.Player = first
	res.State.Queue = append([][]card.Card{second}, res.State.Queue...)
	res.CurrentHandBet = bet
	res.DealerCards = []card.Card{res.State.Dealer}

	return res, nil
}

func finalize(res RoundResolution, engine *Engine) RoundResolution {
	if !res.RoundComplete {
		if len(res.DealerCards) == 0 {
			res.DealerCards = []card.Card{res.State.Dealer}
		}
		return res
	}

	bet := res.CurrentHandBet
	if bet == 0 {
		bet = DefaultBet
	}
	res.State = appendCompleted(res.State, res.Outcome, bet)

	if len(res.State.Queue) > 0 {
		res = advanceQueue(res, bet)
	} else {
		res = settleDealer(res, engine)
	}

	if len(res.DealerCards) == 0 {
		res.DealerCards = []card.Card{res.State.Dealer}
	}

	return res
}

func newResolution(state RoundState) RoundResolution {
	bets := make([]int, len(state.HandBets))
	copy(bets, state.HandBets)
	return RoundResolution{
		State: RoundState{
			Player:    cloneCards(state.Player),
			Dealer:    state.Dealer,
			Queue:     cloneHands(state.Queue),
			Completed: cloneHands(state.Completed),
			Outcomes:  append([]string{}, state.Outcomes...),
			HandBets:  bets,
		},
	}
}

func appendCompleted(state RoundState, outcome string, bet int) RoundState {
	current := cloneCards(state.Player)
	state.Completed = append(state.Completed, current)

	if outcome == "" && hand.NewHand(current).IsBust() {
		outcome = "Bust"
	}

	state.Outcomes = append(state.Outcomes, outcome)
	state.HandBets = append(state.HandBets, bet)
	return state
}

func advanceQueue(res RoundResolution, bet int) RoundResolution {
	next := cloneCards(res.State.Queue[0])
	res.State.Queue = cloneHands(res.State.Queue[1:])
	res.State.Player = next
	res.RoundComplete = false
	res.Outcome = ""
	res.CurrentHandBet = bet
	res.DealerCards = []card.Card{res.State.Dealer}
	return res
}

func settleDealer(res RoundResolution, engine *Engine) RoundResolution {
	dealerCards, outcomes := engine.EvaluateAllHands(res.State.Completed, res.State.Dealer)
	res.DealerCards = dealerCards

	for i := range res.State.Outcomes {
		if res.State.Outcomes[i] == "" && i < len(outcomes) {
			res.State.Outcomes[i] = outcomes[i]
		}
	}

	for i := len(res.State.Outcomes); i < len(outcomes); i++ {
		res.State.Outcomes = append(res.State.Outcomes, outcomes[i])
		if len(res.State.HandBets) <= i {
			res.State.HandBets = append(res.State.HandBets, res.CurrentHandBet)
		}
	}

	res.Outcome = FormatOutcomeSummary(res.State.Outcomes)
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
