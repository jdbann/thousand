package handlers

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	InternalServerError = NewError("Internal Server Error", http.StatusInternalServerError)
)

type HTTPError struct {
	publicMsg string
	status    int
	cause     error
}

func NewError(msg string, status int) HTTPError {
	return HTTPError{
		publicMsg: msg,
		status:    status,
	}
}

func (err HTTPError) Error() string {
	if err.cause != nil {
		return fmt.Sprintf("%s: %s", err.publicMsg, err.cause.Error())
	}

	return err.publicMsg
}

func (err HTTPError) Is(target error) bool {
	var t HTTPError
	if errors.As(target, &t) {
		return err.publicMsg == t.publicMsg && err.status == t.status
	} else {
		return false
	}
}

func (err HTTPError) PublicError() string {
	return fmt.Sprintf("%d: %s", err.status, err.publicMsg)
}

func (err HTTPError) Status() int {
	return err.status
}

func handleError(w http.ResponseWriter, err error) {
	httpErr := HTTPError{}
	if !errors.As(err, &httpErr) {
		httpErr = InternalServerError
	}

	http.Error(w, httpErr.PublicError(), httpErr.Status())
}
