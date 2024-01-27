package dtos

type DeleteBatchBody struct {
	IDs []int `json:"ids"`
}

func (b *DeleteBatchBody) Valid() bool {
	if len(b.IDs) == 0 {
		return false
	}

	return true
}
