package dtos

import "github.com/MowlCoder/goph-keeper/internal/domain"

type AddNewTextBody struct {
	Data domain.TextData `json:"data"`
	Meta string          `json:"meta"`
}

func (b *AddNewTextBody) Valid() bool {
	if b.Data.Text == "" {
		return false
	}

	return true
}

func (b *AddNewTextBody) GetMeta() string {
	return b.Meta
}

func (b *AddNewTextBody) GetData() interface{} {
	return b.Data
}
