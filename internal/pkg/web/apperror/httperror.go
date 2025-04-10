package apperror

import "net/http"

type HTTPError struct {
	Code    int
	Message string
}

func (e HTTPError) Error() string {
	return e.Message
}

func NewInternal(msg string) error {
	return HTTPError{Code: http.StatusInternalServerError, Message: msg}
}

func NewBadReq(msg string) error {
	return HTTPError{Code: http.StatusBadRequest, Message: msg}
}

func NewUnauthorized(msg string) error {
	return HTTPError{Code: http.StatusUnauthorized, Message: msg}
}
