package dtos

import "github.com/MowlCoder/goph-keeper/internal/domain"

type AddNewCardBody struct {
	Data domain.CardData `json:"data"`
	Meta string          `json:"meta"`
}

func (b *AddNewCardBody) Valid() bool {
	if b.Data.Number == "" || b.Data.ExpiredAt == "" || b.Data.CVV == "" {
		return false
	}

	return true
}

func (b *AddNewCardBody) GetMeta() string {
	return b.Meta
}

func (b *AddNewCardBody) GetData() interface{} {
	return b.Data
}
