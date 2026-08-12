package strategy

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func TestMakeDecision_EmptyHands(t *testing.T) {
	tests := []struct {
		name       string
		playerHand *hand.Hand
		dealerHand *hand.Hand
	}{
		{
			"empty player hand",
			hand.NewHand([]card.Card{}),
			hand.NewHand([]card.Card{card.NewCard(card.Ten)}),
		},
		{
			"empty dealer hand",
			hand.NewHand([]card.Card{card.NewCard(card.Ten)}),
			hand.NewHand([]card.Card{}),
		},
		{
			"both empty",
			hand.NewHand([]card.Card{}),
			hand.NewHand([]card.Card{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := MakeDecision(tt.playerHand, tt.dealerHand)
			if err == nil || decision != Hit {
				t.Errorf("MakeDecision() = (%v, %v), want (%v, error)", decision, err, Hit)
			}
		})
	}
}

func TestMakeDecision_Blackjack(t *testing.T) {
	playerHand := hand.NewHand([]card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten)})
	dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

	decision, err := MakeDecision(playerHand, dealerHand)
	if err != nil {
		t.Errorf("MakeDecision() error = %v, want nil", err)
	}
	if decision != Stand {
		t.Errorf("MakeDecision() = %v, want %v", decision, Stand)
	}
}

func TestMakeDecision_Bust(t *testing.T) {
	playerHand := hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten), card.NewCard(card.Five)})
	dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

	decision, err := MakeDecision(playerHand, dealerHand)
	if err != nil {
		t.Errorf("MakeDecision() error = %v, want nil", err)
	}
	if decision != Stand {
		t.Errorf("MakeDecision() = %v, want %v", decision, Stand)
	}
}

func TestMakeDecision_DoubleDownConversion(t *testing.T) {
	tests := []struct {
		name         string
		playerCards  []card.Card
		wantDecision Decision
	}{
		{
			"double down with 2 cards",
			[]card.Card{card.NewCard(card.Seven), card.NewCard(card.Four)},
			DoubleDown,
		},
		{
			"double down with 3 cards converts to hit",
			[]card.Card{card.NewCard(card.Five), card.NewCard(card.Four), card.NewCard(card.Two)},
			Hit,
		},
		{
			"double down with 4 cards converts to hit",
			[]card.Card{card.NewCard(card.Two), card.NewCard(card.Three), card.NewCard(card.Four), card.NewCard(card.Two)},
			Hit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playerHand := hand.NewHand(tt.playerCards)
			dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

			decision, err := MakeDecision(playerHand, dealerHand)
			if err != nil {
				t.Errorf("MakeDecision() error = %v, want nil", err)
			}
			if decision != tt.wantDecision {
				t.Errorf("MakeDecision() = %v, want %v", decision, tt.wantDecision)
			}
		})
	}
}
