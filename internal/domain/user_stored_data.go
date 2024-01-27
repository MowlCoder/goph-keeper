package domain

import "time"

const (
	LogPassDataType = "logpass"
	CardDataType    = "card"
	TextDataType    = "text"
	FileDataType    = "file"
)

type UserStoredData struct {
	ID          int         `json:"id"`
	UserID      int         `json:"user_id"`
	DataType    string      `json:"data_type"`
	Data        interface{} `json:"data"`
	PathOnDisc  string      `json:"path_on_disc,omitempty"`
	CryptedData []byte      `json:"crypted_data,omitempty"`
	Meta        string      `json:"meta"`
	Version     int         `json:"version"`
	CreatedAt   time.Time   `json:"created_at"`
}

func (data UserStoredData) IsLocal() bool {
	return data.Version == -1 || data.ID < 0
}

type LogPassData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type CardData struct {
	Number    string `json:"number"`
	ExpiredAt string `json:"expired_at"`
	CVV       string `json:"cvv"`
}

type TextData struct {
	Text string `json:"text"`
}

type FileData struct {
	Content []byte `json:"content"`
	Name    string `json:"name"`
}

type AddUserStoredDataBody interface {
	Valid() bool
	GetData() interface{}
	GetMeta() string
}
