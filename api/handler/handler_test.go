package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
)

func TestDealHidesDealerHoleCard(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame()
	store.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/deal", nil)
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()
	NewDeal(store)(res, req)
	var body struct {
		DealerCard string `json:"dealerCard"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.DealerCard == "" {
		t.Fatal("dealer upcard missing")
	}
	store.WithGame(g.ID, func(g *game.Game) {
		if len(g.DealerHand.Cards) != 2 {
			t.Fatalf("dealer cards = %d", len(g.DealerHand.Cards))
		}
	})
}

func TestMoveHidesDealerHoleCardWhileRoundActive(t *testing.T) {
	sessions := store.NewSessionStore()
	g := game.NewGame()
	g.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})
	g.UnresolvedHands = []*hand.Hand{
		hand.NewHand([]card.Card{card.NewCard(card.Nine), card.NewCard(card.Eight)}),
	}
	g.DealerHand = hand.NewHand([]card.Card{card.NewCard(card.Six), card.NewCard(card.Ten)})
	g.RoundState = game.RoundStateActive
	sessions.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/move", strings.NewReader(`{"move":"STAND"}`))
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()

	NewMove(sessions)(res, req)

	var body struct {
		RoundState  string         `json:"roundState"`
		DealerCards []string       `json:"dealerCards"`
		Outcomes    []game.Outcome `json:"outcomes"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.RoundState != string(game.RoundStateActive) {
		t.Fatalf("round state = %q, want active", body.RoundState)
	}
	if len(body.DealerCards) != 1 || body.DealerCards[0] != g.DealerHand.Cards[0].ToString() {
		t.Fatalf("dealer cards = %v, want upcard only", body.DealerCards)
	}
	if body.Outcomes == nil || len(body.Outcomes) != 0 {
		t.Fatalf("outcomes = %v, want empty array", body.Outcomes)
	}
}

func TestAbandonCompletesActiveRound(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame()
	g.StartRound()
	store.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/abandon", nil)
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()
	NewAbandon(store)(res, req)
	if res.Code != http.StatusNoContent {
		t.Fatalf("status = %d", res.Code)
	}
}
