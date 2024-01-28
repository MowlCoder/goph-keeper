package domain

import (
	"encoding/json"
	"time"
)

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

func ParseUserStoredData(dataType string, data []byte) (interface{}, error) {
	switch dataType {
	case LogPassDataType:
		var parsedData LogPassData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case CardDataType:
		var parsedData CardData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case TextDataType:
		var parsedData TextData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	case FileDataType:
		var parsedData FileData
		if err := json.Unmarshal(data, &parsedData); err != nil {
			return nil, err
		}

		return parsedData, nil
	default:
		return nil, ErrInvalidDataType
	}
}
