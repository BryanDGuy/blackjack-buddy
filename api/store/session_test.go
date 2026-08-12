package store

import (
	"testing"
	"time"

	"github.com/bryan/blackjack-buddy/internal/dealer"
	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/bryan/blackjack-buddy/internal/player"
)

func TestSessionStoreExpiresIdleGames(t *testing.T) {
	now := time.Date(2026, 8, 12, 0, 0, 0, 0, time.UTC)
	store := newSessionStore(func() time.Time { return now })
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	store.Create(g)
	now = now.Add(24*time.Hour + time.Second)
	if store.WithGame(g.ID, func(*game.Game) {}) {
		t.Fatal("expired game still exists")
	}
}

func TestSessionStoreSerializesOneGame(t *testing.T) {
	store := NewSessionStore()
	g := game.NewGame(player.NewPlayer(), dealer.NewDealer())
	store.Create(g)
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go store.WithGame(g.ID, func(*game.Game) {
		close(entered)
		<-release
	})
	<-entered
	attempted := make(chan struct{})
	go func() {
		close(attempted)
		store.WithGame(g.ID, func(*game.Game) { close(done) })
	}()
	<-attempted
	select {
	case <-done:
		t.Fatal("second operation entered while first held the game")
	default:
	}
	close(release)
	<-done
}
