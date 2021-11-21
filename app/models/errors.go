package models

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
)

var ErrNotFound = NewError("Not found")

type Error struct {
	msg   string
	cause error
}

func NewError(msg string) Error {
	return Error{msg: msg}
}

func (err Error) Cause(cause error) Error { err.cause = cause; return err }

func (err Error) Error() string {
	msg := "Not found"

	if err.cause != nil {
		return fmt.Sprintf("%s: %s", msg, err.cause.Error())
	}

	return msg
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

// ErrMemoryFull is returned when trying to add experiences to a full memory.
var ErrMemoryFull = errors.New("memory is full")

func isMemoryFullError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "TH001"
	}

	return false
}
