package store

import (
	"sync"

	"github.com/bryan/blackjack-buddy/internal/game"
)

type SessionStore struct {
	mu    sync.RWMutex
	games map[string]*game.Game
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		games: make(map[string]*game.Game),
	}
}

func (s *SessionStore) Create(game *game.Game) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[game.ID] = game
	return game.ID
}

func (s *SessionStore) Get(id string) (*game.Game, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	game, exists := s.games[id]
	return game, exists
}
