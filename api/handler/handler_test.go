package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/card"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/hand"
	"github.com/bryan/blackjack-buddy/internal/player"
)

func TestDealHidesDealerHoleCard(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(&player.Player{}, &dealer.Dealer{})
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
		if len(g.Dealer.Hand.Cards) != 2 {
			t.Fatalf("dealer cards = %d", len(g.Dealer.Hand.Cards))
		}
	})
}

func TestMoveHidesDealerHoleCardWhileRoundActive(t *testing.T) {
	sessions := store.NewSessionStore()
	p := &player.Player{}
	p.ActiveHand = hand.NewHand([]card.Card{card.NewCard(card.Ten), card.NewCard(card.Seven)})
	p.UnresolvedHands = []*hand.Hand{
		hand.NewHand([]card.Card{card.NewCard(card.Nine), card.NewCard(card.Eight)}),
	}
	d := &dealer.Dealer{}
	d.Hand = hand.NewHand([]card.Card{card.NewCard(card.Six), card.NewCard(card.Ten)})
	g := game.NewGame(p, d)
	g.RoundState = game.RoundStateActive
	sessions.Create(g)
	req := httptest.NewRequest(http.MethodPost, "/api/game/"+g.ID+"/move", strings.NewReader(`{"move":"STAND"}`))
	req.SetPathValue("id", g.ID)
	res := httptest.NewRecorder()

	NewMove(sessions)(res, req)

	var body struct {
		RoundState  string   `json:"roundState"`
		DealerCards []string `json:"dealerCards"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.RoundState != string(game.RoundStateActive) {
		t.Fatalf("round state = %q, want active", body.RoundState)
	}
	if len(body.DealerCards) != 1 || body.DealerCards[0] != d.Hand.Cards[0].ToString() {
		t.Fatalf("dealer cards = %v, want upcard only", body.DealerCards)
	}
}

func TestAbandonCompletesActiveRound(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(&player.Player{}, &dealer.Dealer{})
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
