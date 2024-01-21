package session

import "sync"

type ClientSession struct {
	token             string
	deletedLogPassIDs map[int]struct{}

	mu *sync.RWMutex
}

func NewClientSession() *ClientSession {
	return &ClientSession{
		deletedLogPassIDs: map[int]struct{}{},

		mu: &sync.RWMutex{},
	}
}

func (s *ClientSession) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.token = token
}

func (s *ClientSession) IsAuth() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.token != ""
}

func (s *ClientSession) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.token
}

func (s *ClientSession) AddDeletedLogPassID(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.deletedLogPassIDs[id] = struct{}{}
}

func (s *ClientSession) IsLogPassDeleted(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.deletedLogPassIDs[id]
	return ok
}

func (s *ClientSession) ClearLogPassDeleted() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.deletedLogPassIDs)
}
