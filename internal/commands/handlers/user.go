package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/session"
	"github.com/MowlCoder/goph-keeper/pkg/input"
)

type UserHandler struct {
	httpClient *http.Client
	session    *session.ClientSession
	userApi    userApi
}

type userApi interface {
	Register(ctx context.Context, email string, password string) (string, error)
	Authorize(ctx context.Context, email string, password string) (string, error)
}

func NewUserHandler(
	httpClient *http.Client,
	session *session.ClientSession,
	userApi userApi,
) *UserHandler {
	return &UserHandler{
		httpClient: httpClient,
		session:    session,
		userApi:    userApi,
	}
}

func (h *UserHandler) Register(args []string) error {
	email := input.GetConsoleInput("Enter email: ", "")
	if email == "" {
		return domain.ErrInvalidInputValue
	}

	password := input.GetConsoleInput("Enter password: ", "")
	if password == "" {
		return domain.ErrInvalidInputValue
	}

	token, err := h.userApi.Register(context.Background(), email, password)
	if err != nil {
		return err
	}

	h.session.SetToken(token)

	return nil
}

func (h *UserHandler) Authorize(args []string) error {
	email := input.GetConsoleInput("Enter email: ", "")
	if email == "" {
		return domain.ErrInvalidInputValue
	}

	password := input.GetConsoleInput("Enter password: ", "")
	if password == "" {
		return domain.ErrInvalidInputValue
	}

	token, err := h.userApi.Authorize(context.Background(), email, password)
	if err != nil {
		return err
	}

	fmt.Println("You successfully authorized.")

	h.session.SetToken(token)

	return nil
}
