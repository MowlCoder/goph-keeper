package domain

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrWrongCredentials  = errors.New("wrong credentials")
	ErrInvalidBody       = errors.New("invalid body")
	ErrNotAuth           = errors.New("not authorized")
	ErrInternal          = errors.New("internal error")

	ErrUserStoredDataNotFound = errors.New("user stored data not found")
	ErrInvalidDataType        = errors.New("invalid data type")

	ErrInvalidCardNumber    = errors.New("invalid card number")
	ErrInvalidCardExpiredAt = errors.New("invalid card expired at (e.g. 4/30)")
	ErrInvalidCardCVV       = errors.New("invalid card cvv")

	ErrCommandNotFound     = errors.New("command not found")
	ErrQuitApp             = errors.New("requested quit from the app")
	ErrInvalidCommandUsage = errors.New("invalid command usage")
	ErrInvalidInputValue   = errors.New("invalid input value")
)
