package lknpd

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
