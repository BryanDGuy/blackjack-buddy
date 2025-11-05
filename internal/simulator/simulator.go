package simulator

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

type Result struct {
	FinalPot        float64
	RoundsPlayed    int
	RanOutOfMoney   bool
	GainLossPercent float64
}

type RoundResult struct {
	playerValue int
	dealerValue int
	playerBust  bool
	dealerBust  bool
	blackjack   bool
}

type Simulator struct {
	strategy strategy.Strategy
	rng      *rand.Rand
}

func NewSimulator(strat strategy.Strategy) *Simulator {
	return &Simulator{
		strategy: strat,
		rng:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *Simulator) Run(startingPot, buyIn float64, rounds int, verbose bool) Result {
	pot := startingPot
	roundsPlayed := 0

	for i := 0; i < rounds; i++ {
		if pot < buyIn {
			return Result{
				FinalPot:        pot,
				RoundsPlayed:    roundsPlayed,
				RanOutOfMoney:   true,
				GainLossPercent: ((pot - startingPot) / startingPot) * 100,
			}
		}

		pot -= buyIn
		outcome, roundResult := s.playRound()
		pot += outcome * buyIn
		roundsPlayed++

		if verbose {
			s.printRound(i+1, roundResult, outcome*buyIn-buyIn, pot)
		}
	}

	return Result{
		FinalPot:        pot,
		RoundsPlayed:    roundsPlayed,
		RanOutOfMoney:   false,
		GainLossPercent: ((pot - startingPot) / startingPot) * 100,
	}
}

func (s *Simulator) playRound() (float64, RoundResult) {
	playerHand := hand.NewHand([]card.Card{s.dealCard(), s.dealCard()})
	dealerHand := hand.NewHand([]card.Card{s.dealCard()})

	if playerHand.IsBlackjack() {
		dealerSecond := s.dealCard()
		fullDealerHand := hand.NewHand(append(dealerHand.Cards, dealerSecond))
		rr := RoundResult{
			playerValue: 21,
			dealerValue: fullDealerHand.Value(),
			blackjack:   true,
		}
		if fullDealerHand.IsBlackjack() {
			return 1.0, rr
		}
		return 2.5, rr
	}

	playerHand = s.playPlayerHand(playerHand, dealerHand)
	if playerHand.IsBust() {
		return 0.0, RoundResult{
			playerValue: playerHand.Value(),
			playerBust:  true,
		}
	}

	dealerFinalHand := s.playDealerHand(dealerHand)
	if dealerFinalHand.IsBust() {
		return 2.0, RoundResult{
			playerValue: playerHand.Value(),
			dealerValue: dealerFinalHand.Value(),
			dealerBust:  true,
		}
	}

	playerValue := playerHand.Value()
	dealerValue := dealerFinalHand.Value()

	rr := RoundResult{
		playerValue: playerValue,
		dealerValue: dealerValue,
	}

	if playerValue > dealerValue {
		return 2.0, rr
	} else if playerValue < dealerValue {
		return 0.0, rr
	}
	return 1.0, rr
}

func (s *Simulator) playPlayerHand(playerHand, dealerHand *hand.Hand) *hand.Hand {
	for !playerHand.IsBust() && !playerHand.IsBlackjack() && playerHand.Value() < 21 {
		decision := s.strategy.GetDecision(playerHand, dealerHand)

		if decision == strategy.Stand {
			break
		}

		if decision == strategy.Hit || decision == strategy.DoubleDown {
			newCards := append(playerHand.Cards, s.dealCard())
			playerHand = hand.NewHand(newCards)

			if decision == strategy.DoubleDown {
				break
			}
		}

		if decision == strategy.Split {
			newCards := append(playerHand.Cards, s.dealCard())
			playerHand = hand.NewHand(newCards)
		}
	}

	return playerHand
}

func (s *Simulator) playDealerHand(dealerHand *hand.Hand) *hand.Hand {
	dealerHand = hand.NewHand(append(dealerHand.Cards, s.dealCard()))

	for dealerHand.Value() < 17 {
		dealerHand = hand.NewHand(append(dealerHand.Cards, s.dealCard()))
	}

	return dealerHand
}

func (s *Simulator) dealCard() card.Card {
	ranks := []card.Rank{
		card.Ace, card.Two, card.Three, card.Four, card.Five,
		card.Six, card.Seven, card.Eight, card.Nine, card.Ten,
		card.Jack, card.Queen, card.King,
	}
	return card.NewCard(ranks[s.rng.Intn(len(ranks))])
}

func (s *Simulator) printRound(roundNum int, rr RoundResult, netChange, pot float64) {
	if rr.blackjack {
		fmt.Printf("%d | BJ vs %d | %+.2f | %.2f\n", roundNum, rr.dealerValue, netChange, pot)
	} else if rr.playerBust {
		fmt.Printf("%d | %d(B) vs ? | %+.2f | %.2f\n", roundNum, rr.playerValue, netChange, pot)
	} else if rr.dealerBust {
		fmt.Printf("%d | %d vs %d(B) | %+.2f | %.2f\n", roundNum, rr.playerValue, rr.dealerValue, netChange, pot)
	} else {
		fmt.Printf("%d | %d vs %d | %+.2f | %.2f\n", roundNum, rr.playerValue, rr.dealerValue, netChange, pot)
	}
}
