package errorz

import (
	"fmt"
	"net/http"
)

var (
	ErrBadRequest              = NewAPIError(http.StatusBadRequest, "Bad Request")
	ErrUnauthorized            = NewAPIError(http.StatusUnauthorized, "Unauthorized")
	ErrForbidden               = NewAPIError(http.StatusForbidden, "Forbidden")
	ErrNotFound                = NewAPIError(http.StatusNotFound, "Not Found")
	ErrConflict                = NewAPIError(http.StatusConflict, "Conflict")
	ErrUnprocessable           = NewAPIError(http.StatusUnprocessableEntity, "Unable to process the contained instructions")
	ErrInternal                = NewAPIError(http.StatusInternalServerError, "Internal Server Error")
	ErrServiceUnavailable      = NewAPIError(http.StatusServiceUnavailable, "Service Unavailable")
	ErrAlreadyExists           = JoinError("Already Exists", ErrConflict)
	ErrSomethingWentWrong      = JoinError("something went wrong", ErrInternal)
	ErrInvalidTransactionState = JoinError("invalid transaction state", ErrInternal)
)

func NewNotFoundError(message string) error {
	return ErrNotFound.JoinError(message)
}

func NewBadRequestError(message string) error {
	return ErrBadRequest.JoinError(message)
}

func NewUnauthorizedError(message string) error {
	return ErrUnauthorized.JoinError(message)
}

func NewForbiddenError(message string) error {
	return ErrForbidden.JoinError(message)
}

func NewInternalError(message string) error {
	return ErrInternal.JoinError(message)
}

func JoinError(message string, base error) error {
	if base.Error() == "" {
		return fmt.Errorf("%v%w", message, base)
	}
	return fmt.Errorf("%v %w", message, base)
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	if message == "" {
		return err
	}
	return fmt.Errorf("%s: %w", message, err)
}
