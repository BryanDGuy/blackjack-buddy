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
	g := game.NewGame(&player.Player{}, &dealer.Dealer{})
	store.Create(g)
	now = now.Add(24*time.Hour + time.Second)
	if store.WithGame(g.ID, func(*game.Game) {}) {
		t.Fatal("expired game still exists")
	}
}

func TestSessionStoreSerializesOneGame(t *testing.T) {
	store := NewSessionStore()
	g := game.NewGame(&player.Player{}, &dealer.Dealer{})
	store.Create(g)
	entry := store.games[g.ID]
	entered := make(chan struct{})
	release := make(chan struct{})
	done := make(chan struct{})
	go func() {
		store.WithGame(g.ID, func(*game.Game) {
			close(entered)
			<-release
		})
		close(done)
	}()
	<-entered
	secondEntered := entry.mu.TryLock()
	if secondEntered {
		entry.mu.Unlock()
	}
	close(release)
	<-done
	if secondEntered {
		t.Fatal("second operation entered while first held the game")
	}
}
