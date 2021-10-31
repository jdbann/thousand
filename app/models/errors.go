package models

import (
	"errors"

	"github.com/jackc/pgconn"
)

// ErrMemoryFull is returned when trying to add experiences to a full memory.
var ErrMemoryFull = errors.New("memory is full")

func isMemoryFullError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "TH001"
	}

	return false
}
