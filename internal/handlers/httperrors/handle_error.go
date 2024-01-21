package httperrors

import (
	"net/http"

	"github.com/MowlCoder/goph-keeper/internal/domain"
	"github.com/MowlCoder/goph-keeper/pkg/httputils"
)

type errorInfo struct {
	statusCode int
	errorCode  int
}

var errorToErrorInfo = map[error]errorInfo{
	domain.ErrWrongCredentials: {
		statusCode: http.StatusBadRequest,
		errorCode:  1,
	},
	domain.ErrUserNotFound: {
		statusCode: http.StatusNotFound,
		errorCode:  2,
	},
	domain.ErrEmailAlreadyTaken: {
		statusCode: http.StatusConflict,
		errorCode:  3,
	},
	domain.ErrInvalidBody: {
		statusCode: http.StatusBadRequest,
		errorCode:  4,
	},
	domain.ErrNotAuth: {
		statusCode: http.StatusUnauthorized,
		errorCode:  401,
	},
	domain.ErrInternal: {
		statusCode: http.StatusInternalServerError,
		errorCode:  500,
	},
}

func Handle(w http.ResponseWriter, err error) {
	var info errorInfo

	if i, ok := errorToErrorInfo[err]; ok {
		info = i
	} else {
		info = errorToErrorInfo[domain.ErrInternal]
	}

	httputils.SendJSONErrorResponse(w, info.statusCode, err.Error(), info.errorCode)
}
