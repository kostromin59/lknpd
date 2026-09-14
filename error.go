package lknpd

import (
	"errors"
	"fmt"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	Code       string
	Message    string
	StatusCode int
}

func (e Error) Error() string {
	return fmt.Sprintf("%s (%d): %s", e.Code, e.StatusCode, e.Message)
}
