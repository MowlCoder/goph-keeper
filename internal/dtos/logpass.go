package dtos

type AddNewLogPassBody struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Source   string `json:"source"`
}

func (b *AddNewLogPassBody) Valid() bool {
	if b.Login == "" || b.Password == "" {
		return false
	}

	return true
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
