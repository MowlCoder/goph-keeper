package session

import "sync"

type ClientSession struct {
	mu *sync.RWMutex

	token      string
	deletedIDs map[int]struct{}
	editedIDs  map[int]struct{}
}

func NewClientSession() *ClientSession {
	return &ClientSession{
		mu: &sync.RWMutex{},

		deletedIDs: map[int]struct{}{},
		editedIDs:  map[int]struct{}{},
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

func (s *ClientSession) AddDeleted(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.deletedIDs[id] = struct{}{}
}

func (s *ClientSession) IsDeleted(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.deletedIDs[id]
	return ok
}

func (s *ClientSession) ClearDeleted() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.deletedIDs)
}

func (s *ClientSession) AddEdited(id int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.editedIDs[id] = struct{}{}
}

func (s *ClientSession) IsEdited(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.editedIDs[id]
	return ok
}

func (s *ClientSession) ClearEdited() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.editedIDs)
}
