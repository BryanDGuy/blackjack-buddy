package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/strategy"
)

func main() {
	fmt.Println("🃏 Welcome to Blackjack Buddy!")
	fmt.Println("A Blackjack decision assistant to help you make optimal plays.")
	fmt.Println()

	advisor := strategy.NewAdvisor()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("=== New Hand ===")

		playerHand := getPlayerHand(scanner)
		if playerHand == nil {
			fmt.Println("Goodbye!")
			return
		}

		dealerHand := getDealerHand(scanner)
		if dealerHand == nil {
			fmt.Println("Goodbye!")
			return
		}

		decision, err := advisor.GetDecision(playerHand, dealerHand)
		if err != nil {
			fmt.Printf("Error getting advice: %v\n", err)
			continue
		}

		fmt.Println()
		fmt.Printf("Your hand: %s (Value: %d)\n", playerHand.ToString(), playerHand.Value())
		fmt.Printf("Dealer's up card: %s\n", dealerHand.Cards[0].ToString())
		fmt.Println()
		fmt.Printf("🎯 RECOMMENDED ACTION: %s\n", decision.ToString())

		provideContext(decision, playerHand)

		fmt.Println()
		fmt.Println("Press Enter to analyze another hand, or type 'quit' to exit...")
		if !scanner.Scan() {
			break
		}

		if strings.ToLower(strings.TrimSpace(scanner.Text())) == "quit" {
			fmt.Println("Goodbye!")
			break
		}
	}
}

func getPlayerHand(scanner *bufio.Scanner) *hand.Hand {
	for {
		fmt.Print("Enter your cards (e.g., 'A 7' or '2 3'): ")
		if !scanner.Scan() {
			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" || strings.ToLower(input) == "quit" {
			return nil
		}

		cards, err := parseCards(input)
		if err != nil {
			fmt.Printf("Invalid input: %v\n", err)
			fmt.Println("Please enter cards in format like 'A 7' or 'K 5 3'")
			continue
		}

		playerHand := hand.NewHand(cards)
		if playerHand.IsBust() {
			fmt.Println("Your hand is already bust! Please enter a valid hand.")
			continue
		}

		return playerHand
	}
}

func getDealerHand(scanner *bufio.Scanner) *hand.Hand {
	for {
		fmt.Print("Enter dealer's up card (e.g., 'K' or '10'): ")
		if !scanner.Scan() {
			return nil
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" || strings.ToLower(input) == "quit" {
			return nil
		}

		cards, err := parseCards(input)
		if err != nil || len(cards) != 1 {
			fmt.Printf("Invalid input: %v\n", err)
			fmt.Println("Please enter exactly one card like 'K' or '10'")
			continue
		}

		dealerHand := hand.NewHand(cards)
		return dealerHand
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

func provideContext(decision strategy.Decision, playerHand *hand.Hand) {
	fmt.Println()
	switch decision {
	case strategy.Hit:
		fmt.Println("💡 Take another card. Your current hand isn't strong enough to stand.")
	case strategy.Stand:
		fmt.Println("💡 Stay with your current hand. The risk of busting is too high.")
	case strategy.DoubleDown:
		fmt.Println("💡 Double your bet and take exactly one more card. This is a favorable situation!")
	case strategy.Split:
		fmt.Println("💡 Split your pair into two separate hands. This maximizes your advantage.")
	case strategy.Surrender:
		fmt.Println("💡 Surrender half your bet. This hand has a poor expected value.")
	}

	if playerHand.IsSoft() {
		fmt.Println("📝 Note: You have a 'soft' hand (ace counted as 11).")
	}

	if playerHand.CanSplit() && decision != strategy.Split {
		fmt.Printf("📝 Note: You could split this pair, but it's not recommended.\n")
	}

	if playerHand.IsBlackjack() {
		fmt.Println("🎉 Congratulations! You have a blackjack!")
	}
}
