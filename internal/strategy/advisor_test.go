package strategy

import (
	"testing"

	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

type mockStrategy struct {
	decision Decision
}

func (m *mockStrategy) GetDecision(playerHand, dealerHand *hand.Hand) Decision {
	return m.decision
}

func TestAdvisor_MakeDecision_EmptyHands(t *testing.T) {
	advisor := NewAdvisor(&mockStrategy{decision: Hit})

	tests := []struct {
		name         string
		playerHand   *hand.Hand
		dealerHand   *hand.Hand
		wantError    bool
		wantDecision Decision
	}{
		{
			"empty player hand",
			hand.NewHand([]card.Card{}),
			hand.NewHand([]card.Card{card.NewCard(card.Ten)}),
			true,
			Hit,
		},
		{
			"empty dealer hand",
			hand.NewHand([]card.Card{card.NewCard(card.Ten)}),
			hand.NewHand([]card.Card{}),
			true,
			Hit,
		},
		{
			"both empty",
			hand.NewHand([]card.Card{}),
			hand.NewHand([]card.Card{}),
			true,
			Hit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			decision, err := advisor.MakeDecision(tt.playerHand, tt.dealerHand)
			if (err != nil) != tt.wantError {
				t.Errorf("Advisor.MakeDecision() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if decision != tt.wantDecision {
				t.Errorf("Advisor.MakeDecision() = %v, want %v", decision, tt.wantDecision)
			}
		})
	}
}

func TestAdvisor_MakeDecision_Blackjack(t *testing.T) {
	advisor := NewAdvisor(&mockStrategy{decision: Hit})

	playerHand := hand.NewHand([]card.Card{card.NewCard(card.Ace), card.NewCard(card.Ten)})
	dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

	decision, err := advisor.MakeDecision(playerHand, dealerHand)
	if err != nil {
		t.Errorf("Advisor.MakeDecision() error = %v, want nil", err)
	}
	if decision != Stand {
		t.Errorf("Advisor.MakeDecision() = %v, want %v", decision, Stand)
	}
}

func TestAdvisor_MakeDecision_Bust(t *testing.T) {
	advisor := NewAdvisor(&mockStrategy{decision: Hit})

	playerHand := hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Ten), card.NewCard(card.Five)})
	dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

	decision, err := advisor.MakeDecision(playerHand, dealerHand)
	if err != nil {
		t.Errorf("Advisor.MakeDecision() error = %v, want nil", err)
	}
	if decision != Stand {
		t.Errorf("Advisor.MakeDecision() = %v, want %v", decision, Stand)
	}
}

func TestAdvisor_MakeDecision_DoubleDownConversion(t *testing.T) {
	advisor := NewAdvisor(&mockStrategy{decision: DoubleDown})

	tests := []struct {
		name         string
		playerCards  []card.Card
		wantDecision Decision
	}{
		{
			"double down with 2 cards",
			[]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			DoubleDown,
		},
		{
			"double down with 3 cards converts to hit",
			[]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two)},
			Hit,
		},
		{
			"double down with 4 cards converts to hit",
			[]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven), card.NewCard(card.Two), card.NewCard(card.Two)},
			Hit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			playerHand := hand.NewHand(tt.playerCards)
			dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

			decision, err := advisor.MakeDecision(playerHand, dealerHand)
			if err != nil {
				t.Errorf("Advisor.MakeDecision() error = %v, want nil", err)
			}
			if decision != tt.wantDecision {
				t.Errorf("Advisor.MakeDecision() = %v, want %v", decision, tt.wantDecision)
			}
		})
	}
}

func TestAdvisor_MakeDecision_NormalFlow(t *testing.T) {
	tests := []struct {
		name             string
		strategyDecision Decision
		playerCards      []card.Card
		wantDecision     Decision
	}{
		{
			"strategy returns hit",
			Hit,
			[]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			Hit,
		},
		{
			"strategy returns stand",
			Stand,
			[]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)},
			Stand,
		},
		{
			"strategy returns split",
			Split,
			[]card.Card{card.NewCard(card.Eight), card.NewCard(card.Eight)},
			Split,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advisor := NewAdvisor(&mockStrategy{decision: tt.strategyDecision})
			playerHand := hand.NewHand(tt.playerCards)
			dealerHand := hand.NewHand([]card.Card{card.NewCard(card.Seven)})

			decision, err := advisor.MakeDecision(playerHand, dealerHand)
			if err != nil {
				t.Errorf("Advisor.MakeDecision() error = %v, want nil", err)
			}
			if decision != tt.wantDecision {
				t.Errorf("Advisor.MakeDecision() = %v, want %v", decision, tt.wantDecision)
			}
		})
	}
}
