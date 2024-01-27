package dtos

import "github.com/MowlCoder/goph-keeper/internal/domain"

type AddNewLogPassBody struct {
	Data domain.LogPassData `json:"data"`
	Meta string             `json:"meta"`
}

func (b *AddNewLogPassBody) Valid() bool {
	if b.Data.Login == "" || b.Data.Password == "" {
		return false
	}

	return true
}

func (b *AddNewLogPassBody) GetMeta() string {
	return b.Meta
}

func (b *AddNewLogPassBody) GetData() interface{} {
	return b.Data
}
