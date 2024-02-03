package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/MowlCoder/goph-keeper/internal/dtos"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
)

// UserAPI - struct responsible for communicating with external API
type UserAPI struct {
	baseHTTPAddress string
	httpClient      *http.Client
	session         *session.ClientSession
}

// NewUserAPI - constructor for UserAPI struct
func NewUserAPI(
	baseHTTPAddress string,
	httpClient *http.Client,
	session *session.ClientSession,
) *UserAPI {
	return &UserAPI{
		baseHTTPAddress: baseHTTPAddress,
		httpClient:      httpClient,
		session:         session,
	}
}

// Register - register user and return JWT token
func (api *UserAPI) Register(ctx context.Context, email string, password string) (string, error) {
	body := dtos.RegisterBody{
		Email:    email,
		Password: password,
	}

	if !body.Validate() {
		return "", errors.New("invalid arguments")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := api.httpClient.Post(
		fmt.Sprintf("%s/api/v1/user/register", api.baseHTTPAddress),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return "", err
		}
		return "", errors.New(errResp.Error)
	}

	var respBody dtos.RegisterResponse
	if err := json.Unmarshal(data, &respBody); err != nil {
		return "", err
	}

	return respBody.Token, nil
}

// Authorize - authorize user and return JWT token
func (api *UserAPI) Authorize(ctx context.Context, email string, password string) (string, error) {
	body := dtos.AuthorizeBody{
		Email:    email,
		Password: password,
	}

	if !body.Validate() {
		return "", errors.New("invalid arguments")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	resp, err := api.httpClient.Post(
		fmt.Sprintf("%s/api/v1/user/authorize", api.baseHTTPAddress),
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return "", err
		}
		return "", errors.New(errResp.Error)
	}

	var respBody dtos.AuthorizeResponse
	if err := json.Unmarshal(data, &respBody); err != nil {
		return "", err
	}

	return respBody.Token, nil
}
