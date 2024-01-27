package dtos

type AddNewCardBody struct {
	Number    string `json:"number"`
	ExpiredAt string `json:"expired_at"`
	CVV       string `json:"cvv"`
	Meta      string `json:"meta"`
}

func (b *AddNewCardBody) Valid() bool {
	if b.Number == "" || b.ExpiredAt == "" || b.CVV == "" {
		return false
	}

	return true
}

func (b *AddNewCardBody) GetMeta() string {
	return b.Meta
}

type DeleteBatchBody struct {
	IDs []int `json:"ids"`
}

func (b *DeleteBatchBody) Valid() bool {
	if len(b.IDs) == 0 {
		return false
	}

	return true
}

type DeleteBatchCardsBody struct {
	IDs []int `json:"ids"`
}

func (b *DeleteBatchCardsBody) Valid() bool {
	if len(b.IDs) == 0 {
		return false
	}

	return true
}
