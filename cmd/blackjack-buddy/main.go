package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
	"github.com/bryan/blackjack-buddy/internal/strategy/strategies"
	"github.com/bryan/blackjack-buddy/internal/trainer"
)

func main() {
	trainMode := flag.Bool("train", false, "")
	strategyName := flag.String("strategy", "basic", "")
	flag.Parse()

	stratType := strategies.StrategyType(*strategyName)
	strat := strategies.CreateStrategy(stratType)

	if *trainMode {
		runTrainer(strat)
	} else {
		runInteractive()
	}
}

func runTrainer(strat strategy.Strategy) {
	t := trainer.NewTrainer(strat)
	if err := t.Start(8080); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runInteractive() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Strategy (basic/coward): ")
	if !scanner.Scan() {
		return
	}

	strategyInput := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if strategyInput == "" || strategyInput == "quit" {
		return
	}

	stratType := strategies.StrategyType(strategyInput)
	strat := strategies.CreateStrategy(stratType)
	advisor := strategy.NewAdvisor(strat)

	for {
		playerHand := getPlayerHand(scanner)
		if playerHand == nil {
			return
		}

		dealerHand := getDealerHand(scanner)
		if dealerHand == nil {
			return
		}

		decision, err := advisor.MakeDecision(playerHand, dealerHand)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("%s | Dealer: %s | Action: %s\n",
			playerHand.ToString(),
			dealerHand.Cards[0].ToString(),
			decision.ToString())

		if !scanner.Scan() {
			break
		}
		if strings.ToLower(strings.TrimSpace(scanner.Text())) == "quit" {
			break
		}
	}
}

func getPlayerHand(scanner *bufio.Scanner) *hand.Hand {
	for {
		fmt.Print("Your cards: ")
		if !scanner.Scan() {
			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" || strings.ToLower(input) == "quit" {
			return nil
		}

		cards, err := parseCards(input)
		if err != nil {
			fmt.Printf("Invalid: %v\n", err)
			continue
		}

		playerHand := hand.NewHand(cards)
		if playerHand.IsBust() {
			fmt.Println("Hand is bust")
			continue
		}

		return playerHand
	}
}

func getDealerHand(scanner *bufio.Scanner) *hand.Hand {
	for {
		fmt.Print("Dealer card: ")
		if !scanner.Scan() {
			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" || strings.ToLower(input) == "quit" {
			return nil
		}

		cards, err := parseCards(input)
		if err != nil || len(cards) != 1 {
			fmt.Println("Enter one card")
			continue
		}

		return hand.NewHand(cards)
	}
}

func parseCards(input string) ([]card.Card, error) {
	parts := strings.Fields(strings.ToUpper(input))
	cards := make([]card.Card, 0, len(parts))

	for _, part := range parts {
		var rank card.Rank
		switch part {
		case "A":
			rank = card.Ace
		case "2":
			rank = card.Two
		case "3":
			rank = card.Three
		case "4":
			rank = card.Four
		case "5":
			rank = card.Five
		case "6":
			rank = card.Six
		case "7":
			rank = card.Seven
		case "8":
			rank = card.Eight
		case "9":
			rank = card.Nine
		case "10":
			rank = card.Ten
		case "J":
			rank = card.Jack
		case "Q":
			rank = card.Queen
		case "K":
			rank = card.King
		default:
			return nil, fmt.Errorf("invalid rank: %s", part)
		}

		cards = append(cards, card.NewCard(rank))
	}

	return cards, nil
}
