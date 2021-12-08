package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	msg   string
	cause error
}

func New(msg string) Error {
	return Error{msg: msg}
}

func (err Error) Cause(cause error) Error { err.cause = cause; return err }

func (err Error) Error() string {
	if err.cause != nil {
		return fmt.Sprintf("%s: %s", err.msg, err.cause.Error())
	}

	return err.msg
}

func (err Error) Is(target error) bool {
	var t Error
	if errors.As(target, &t) {
		return err.msg == t.msg
	} else {
		return false
	}
}

func (err Error) Unwrap() error { return err.cause }
