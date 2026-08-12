package store

import (
	"sync"
	"time"

	"github.com/bryan/blackjack-buddy/internal/game"
)

const sessionTTL = 24 * time.Hour

type session struct {
	mu         sync.Mutex
	game       *game.Game
	lastAccess time.Time
}

type SessionStore struct {
	mu    sync.Mutex
	games map[string]*session
	now   func() time.Time
}

func NewSessionStore() *SessionStore {
	return newSessionStore(time.Now)
}

func newSessionStore(now func() time.Time) *SessionStore {
	return &SessionStore{
		games: make(map[string]*session),
		now:   now,
	}
}

func (s *SessionStore) Create(game *game.Game) string {
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()

	s.prune(now)
	s.games[game.ID] = &session{game: game, lastAccess: now}
	return game.ID
}

func (s *SessionStore) WithGame(id string, fn func(*game.Game)) bool {
	now := s.now()
	s.mu.Lock()
	s.prune(now)
	entry, ok := s.games[id]
	if ok {
		entry.lastAccess = now
	}
	s.mu.Unlock()
	if !ok {
		return false
	}

	entry.mu.Lock()
	defer entry.mu.Unlock()
	fn(entry.game)
	return true
}

func (s *SessionStore) prune(now time.Time) {
	for id, entry := range s.games {
		if now.Sub(entry.lastAccess) > sessionTTL {
			delete(s.games, id)
		}
	}
}
