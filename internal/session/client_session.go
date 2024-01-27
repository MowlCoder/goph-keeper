package session

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

type ClientSession struct {
	mu   *sync.RWMutex
	file *os.File

	Token      string           `json:"token"`
	DeletedIDs map[int]struct{} `json:"deleted_ids"`
	EditedIDs  map[int]struct{} `json:"edited_ids"`
}

func NewClientSession(file *os.File) *ClientSession {
	session := &ClientSession{
		mu: &sync.RWMutex{},

		DeletedIDs: map[int]struct{}{},
		EditedIDs:  map[int]struct{}{},
	}

	session.file = file

	if err := json.NewDecoder(file).Decode(&session); err != nil {
		if !errors.Is(err, io.EOF) {
			panic(err)
		}
	}

	session.Token = "" // TODO: validate token

	return session
}

func (s *ClientSession) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Token = token
}

func (s *ClientSession) IsAuth() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Token != ""
}

func (s *ClientSession) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Token
}

func (s *ClientSession) AddDeleted(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DeletedIDs[id] = struct{}{}
	return s.SaveInFile()
}

func (s *ClientSession) IsDeleted(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.DeletedIDs[id]
	return ok
}

func (s *ClientSession) ClearDeleted() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.DeletedIDs)
	return s.SaveInFile()
}

func (s *ClientSession) AddEdited(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.EditedIDs[id] = struct{}{}
	return s.SaveInFile()
}

func (s *ClientSession) IsEdited(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.EditedIDs[id]
	return ok
}

func (s *ClientSession) ClearEdited() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.EditedIDs)
	return s.SaveInFile()
}

func (s *ClientSession) SaveInFile() error {
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	writer := json.NewEncoder(s.file)
	return writer.Encode(s)
}
