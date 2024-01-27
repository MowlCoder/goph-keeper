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

type UserStoredDataAPI struct {
	baseHTTPAddress string
	httpClient      *http.Client
	session         *session.ClientSession
}

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

	return respBody, nil
}

func (api *UserStoredDataAPI) Add(ctx context.Context, entity domain.UserStoredData) (*domain.UserStoredData, error) {
	body := entity.Data
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

func (api *UserStoredDataAPI) DeleteBatch(ctx context.Context, ids []int) error {
	body := &dtos.DeleteBatchCardsBody{
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
