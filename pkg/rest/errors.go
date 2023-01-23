package rest

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrInternalServer   = errors.New("internal server error")
	ErrNotFound         = errors.New("requested content not found")
	ErrStatusBadGateway = errors.New("status bad gateway")
)

type HTTPError struct {
	Code    int
	Message string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("code: %d, error: %s", err.Code, err.Message)
}

func NewInternalServerError() *HTTPError {
	return &HTTPError{
		Code:    http.StatusInternalServerError,
		Message: ErrInternalServer.Error(),
	}
}

func NewNotFoundError() *HTTPError {
	return &HTTPError{
		Code:    http.StatusNotFound,
		Message: ErrNotFound.Error(),
	}
}

func NewBadRequest(msg string) *HTTPError {
	return &HTTPError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
}

func NewStatusBadGateway() *HTTPError {
	return &HTTPError{
		Code:    http.StatusBadGateway,
		Message: ErrStatusBadGateway.Error(),
	}
}
