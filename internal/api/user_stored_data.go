package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/dtos"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
)

// UserStoredDataAPI - struct responsible for communicating with external API
type UserStoredDataAPI struct {
	baseHTTPAddress string
	httpClient      *http.Client
	session         *session.ClientSession
}

// NewUserStoredDataAPI - constructor for UserStoredDataAPI struct
func NewUserStoredDataAPI(
	baseHTTPAddress string,
	httpClient *http.Client,
	session *session.ClientSession,
) *UserStoredDataAPI {
	return &UserStoredDataAPI{
		baseHTTPAddress: baseHTTPAddress,
		httpClient:      httpClient,
		session:         session,
	}
}

// GetAll - get all users records
func (api *UserStoredDataAPI) GetAll(ctx context.Context) ([]domain.UserStoredData, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/data", api.baseHTTPAddress), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+api.session.GetToken())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return nil, err
		}
		return nil, errors.New(errResp.Error)
	}

	var respBody []domain.UserStoredData
	if err := json.Unmarshal(data, &respBody); err != nil {
		return nil, err
	}

	for idx, data := range respBody {
		respBody[idx].Data, err = api.parseData(data)
		if err != nil {
			return nil, err
		}
	}

	return respBody, nil
}

type addBody struct {
	Data interface{} `json:"data"`
	Meta string      `json:"meta"`
}

// Add - add user record to external service
func (api *UserStoredDataAPI) Add(ctx context.Context, entity domain.UserStoredData) (*domain.UserStoredData, error) {
	body := &addBody{
		Data: entity.Data,
		Meta: entity.Meta,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/data/%s", api.baseHTTPAddress, entity.DataType), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.session.GetToken())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return nil, err
		}
		return nil, errors.New("api error: " + errResp.Error)
	}

	var respBody domain.UserStoredData
	if err := json.Unmarshal(data, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

// UpdateByID - update record with given id at external service
func (api *UserStoredDataAPI) UpdateByID(ctx context.Context, id int, data interface{}, meta string) (*domain.UserStoredData, error) {
	body := &addBody{
		Data: data,
		Meta: meta,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/api/v1/data/update/%d", api.baseHTTPAddress, id), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.session.GetToken())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(respData, &errResp); err != nil {
			return nil, err
		}
		return nil, errors.New("api error: " + errResp.Error)
	}

	var respBody domain.UserStoredData
	if err := json.Unmarshal(respData, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

// DeleteBatch - delete several records with given ids at external service
func (api *UserStoredDataAPI) DeleteBatch(ctx context.Context, ids []int) error {
	body := &dtos.DeleteBatchBody{
		IDs: ids,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/data", api.baseHTTPAddress), bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+api.session.GetToken())

	resp, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return err
		}
		return errors.New(errResp.Error)
	}

	return nil
}

func (api *UserStoredDataAPI) parseData(userData domain.UserStoredData) (interface{}, error) {
	jsonData, err := json.Marshal(userData.Data)
	if err != nil {
		return nil, err
	}

	return domain.ParseUserStoredData(userData.DataType, jsonData)
}
