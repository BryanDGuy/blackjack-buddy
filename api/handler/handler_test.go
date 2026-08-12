package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bryan/blackjack-buddy/api/store"
	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/player"
)

func TestDealHidesDealerHoleCard(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
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

func TestAbandonCompletesActiveRound(t *testing.T) {
	store := store.NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
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
