package handlers

import (
	"bytes"
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

type UserHandler struct {
	httpClient *http.Client
	session    *session.ClientSession
}

func NewUserHandler(
	httpClient *http.Client,
	session *session.ClientSession,
) *UserHandler {
	return &UserHandler{
		httpClient: httpClient,
		session:    session,
	}
}

func (h *UserHandler) Register(args []string) error {
	if len(args) != 2 {
		return domain.ErrInvalidCommandUsage
	}

	body := dtos.RegisterBody{
		Email:    args[0],
		Password: args[1],
	}

	if !body.Validate() {
		return errors.New("invalid arguments")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := h.httpClient.Post(
		"http://localhost:4000/api/v1/user/register",
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return err
		}
		return errors.New(errResp.Error)
	}

	var respBody dtos.RegisterResponse
	if err := json.Unmarshal(data, &respBody); err != nil {
		return err
	}

	fmt.Println(respBody)

	h.session.SetToken(respBody.Token)

	return nil
}

func (h *UserHandler) Authorize(args []string) error {
	if len(args) != 2 {
		return domain.ErrInvalidCommandUsage
	}

	body := dtos.AuthorizeBody{
		Email:    args[0],
		Password: args[1],
	}

	if !body.Validate() {
		return errors.New("invalid arguments")
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	resp, err := h.httpClient.Post(
		"http://localhost:4000/api/v1/user/authorize",
		"application/json",
		bytes.NewReader(b),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var errResp httputils.HTTPError
		if err := json.Unmarshal(data, &errResp); err != nil {
			return err
		}
		return errors.New(errResp.Error)
	}

	var respBody dtos.AuthorizeResponse
	if err := json.Unmarshal(data, &respBody); err != nil {
		return err
	}

	fmt.Println("You successfully authorized.")

	h.session.SetToken(respBody.Token)

	return nil
}
