package dtos

import "github.com/MowlCoder/goph-keeper/internal/domain"

type AddNewFileBody struct {
	Data domain.FileData `json:"data"`
	Meta string          `json:"meta"`
}

func (b *AddNewFileBody) Valid() bool {
	if len(b.Data.Content) == 0 || b.Data.Name == "" {
		return false
	}

	return true
}

func (b *AddNewFileBody) GetMeta() string {
	return b.Meta
}

func (b *AddNewFileBody) GetData() interface{} {
	return b.Data
}
