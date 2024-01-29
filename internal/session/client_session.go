package session

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
)

// ClientSession - struct responsible for keeping client session state
type ClientSession struct {
	mu   *sync.RWMutex
	file *os.File

	Token      string           `json:"token"`
	DeletedIDs map[int]struct{} `json:"deleted_ids"`
	EditedIDs  map[int]struct{} `json:"edited_ids"`
}

// NewClientSession - constructor for ClientSession struct
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

// SetToken - save user token in session state
func (s *ClientSession) SetToken(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Token = token
}

// IsAuth - check if user already authorized
func (s *ClientSession) IsAuth() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Token != ""
}

// GetToken - get user token from session state
func (s *ClientSession) GetToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.Token
}

// AddDeleted - add id of deleted record in session state
func (s *ClientSession) AddDeleted(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.DeletedIDs[id] = struct{}{}
	return s.SaveInFile()
}

// IsDeleted - check if record with given id was deleted
func (s *ClientSession) IsDeleted(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.DeletedIDs[id]
	return ok
}

// ClearDeleted - clear deleted record ids from session state
func (s *ClientSession) ClearDeleted() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.DeletedIDs)
	return s.SaveInFile()
}

// AddEdited - add id of edited record in session state
func (s *ClientSession) AddEdited(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.EditedIDs[id] = struct{}{}
	return s.SaveInFile()
}

// IsEdited - check if record with given id was edited
func (s *ClientSession) IsEdited(id int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.EditedIDs[id]
	return ok
}

// ClearEdited - clear edited record ids from session state
func (s *ClientSession) ClearEdited() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.EditedIDs)
	return s.SaveInFile()
}

// SaveInFile - save session state in file
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
