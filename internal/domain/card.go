package domain

import "time"

type Card struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Number    string    `json:"number"`
	ExpiredAt string    `json:"expired_at"`
	CVV       string    `json:"cvv"`
	CreatedAt time.Time `json:"created_at"`
	Meta      string    `json:"meta"`
	Version   int       `json:"version"`
}

func (c Card) GetID() int {
	return c.ID
}

func (c Card) GetVersion() int {
	return c.Version
}

func (c Card) IsLocal() bool {
	return c.Version == -1 || c.ID < 0
}
