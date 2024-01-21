package domain

import "time"

type LogPass struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	Source    string    `json:"source"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"version"`
}

func (lp LogPass) GetID() int {
	return lp.ID
}

func (lp LogPass) GetVersion() int {
	return lp.Version
}

func (lp LogPass) IsLocal() bool {
	return lp.Version == -1 || lp.ID < 0
}
