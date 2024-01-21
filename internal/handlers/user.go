package handlers

import (
	"context"
	"net/http"

	"github.com/MowlCoder/goph-keeper/internal/dtos"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/internal/handlers/httperrors"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
	jsonutil "github.com/MowlCoder/goph-keeper/pkg/jsonutils"
)

type userService interface {
	Create(ctx context.Context, email string, password string) (*domain.User, error)
	Authorize(ctx context.Context, email string, password string) (*domain.User, error)
}

type tokenGenerator interface {
	Generate(ctx context.Context, user domain.User) (string, error)
}

type UserHandler struct {
	userService    userService
	tokenGenerator tokenGenerator
}

func NewUserHandler(
	userService userService,
	tokenGenerator tokenGenerator,
) *UserHandler {
	return &UserHandler{
		userService:    userService,
		tokenGenerator: tokenGenerator,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body dtos.RegisterBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Validate() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	user, err := h.userService.Create(r.Context(), body.Email, body.Password)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	token, err := h.tokenGenerator.Generate(r.Context(), *user)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusCreated, dtos.RegisterResponse{
		Token: token,
	})
}

func (h *UserHandler) Authorize(w http.ResponseWriter, r *http.Request) {
	var body dtos.AuthorizeBody
	if statusCode, err := jsonutil.Unmarshal(w, r, &body); err != nil {
		httputils.SendJSONErrorResponse(w, statusCode, err.Error(), statusCode)
		return
	}

	if !body.Validate() {
		httperrors.Handle(w, domain.ErrInvalidBody)
		return
	}

	user, err := h.userService.Authorize(r.Context(), body.Email, body.Password)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	token, err := h.tokenGenerator.Generate(r.Context(), *user)
	if err != nil {
		httperrors.Handle(w, err)
		return
	}

	httputils.SendJSONResponse(w, http.StatusOK, dtos.AuthorizeResponse{
		Token: token,
	})
}
