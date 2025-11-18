package store

import (
	"sync"

	"github.com/bryan/blackjack-buddy/internal/game"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*game.Session
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*game.Session),
	}
}

func (s *SessionStore) Create(session *game.Session) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.sessions[session.ID()] = session
	return session.ID()
}

func (s *SessionStore) Get(id string) (*game.Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	return session, exists
}
