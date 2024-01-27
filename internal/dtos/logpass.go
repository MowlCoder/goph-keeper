package dtos

type AddNewLogPassBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Meta     string `json:"meta"`
}

func (b *AddNewLogPassBody) Valid() bool {
	if b.Login == "" || b.Password == "" {
		return false
	}

	return true
}

func (b *AddNewLogPassBody) GetMeta() string {
	return b.Meta
}

type DeleteBatchPairsBody struct {
	IDs []int `json:"ids"`
}

func (b *DeleteBatchPairsBody) Valid() bool {
	if len(b.IDs) == 0 {
		return false
	}

	return true
}
