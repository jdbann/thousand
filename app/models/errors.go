package models

import "emailaddress.horse/thousand/errors"

var (
	PgErrCodeMemoryFull = "TH001"
)

var (
	// ErrNotFound is returned when a requested record could not be found.
	ErrNotFound = errors.New("Not found")

	// ErrMemoryFull is returned when trying to add experiences to a full memory.
	ErrMemoryFull = errors.New("Memory is full")
)
