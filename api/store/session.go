package store

import (
	"sync"

	"github.com/bryan/blackjack-buddy/internal/game"
	"github.com/google/uuid"
)

type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*game.GameSession
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*game.GameSession),
	}
}

func (s *SessionStore) Create(session *game.GameSession) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	s.sessions[id] = session
	return id
}

func (s *SessionStore) Get(id string) (*game.GameSession, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.sessions[id]
	return session, exists
}

