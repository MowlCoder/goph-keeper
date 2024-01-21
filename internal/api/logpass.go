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

type LogPassAPI struct {
	baseHTTPAddress string
	httpClient      *http.Client
	session         *session.ClientSession
}

func NewLogPassAPI(
	baseHTTPAddress string,
	httpClient *http.Client,
	session *session.ClientSession,
) *LogPassAPI {
	return &LogPassAPI{
		baseHTTPAddress: baseHTTPAddress,
		httpClient:      httpClient,
		session:         session,
	}
}

type getUserPairsResponse struct {
	Data []domain.LogPass
}

func (api *LogPassAPI) GetUserPairs(ctx context.Context) ([]domain.LogPass, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/v1/logpass", api.baseHTTPAddress), nil)
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

	var respBody getUserPairsResponse
	if err := json.Unmarshal(data, &respBody); err != nil {
		return nil, err
	}

	return respBody.Data, nil
}

func (api *LogPassAPI) AddNewPair(
	ctx context.Context,
	login string,
	password string,
	source string,
) (*domain.LogPass, error) {
	body := &dtos.AddNewLogPassBody{
		Login:    login,
		Password: password,
		Source:   source,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/logpass", api.baseHTTPAddress), bytes.NewReader(b))
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
		return nil, errors.New(errResp.Error)
	}

	var respBody domain.LogPass
	if err := json.Unmarshal(data, &respBody); err != nil {
		return nil, err
	}

	return &respBody, nil
}

func (api *LogPassAPI) DeleteBatchPairs(
	ctx context.Context,
	ids []int,
) error {
	body := &dtos.DeleteBatchPairsBody{
		IDs: ids,
	}
	b, _ := json.Marshal(body)

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api/v1/logpass", api.baseHTTPAddress), bytes.NewReader(b))
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
