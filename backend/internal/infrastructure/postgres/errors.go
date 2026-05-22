package postgres

import (
	"database/sql"
	"errors"
)

var (
	// ErrNotFound is returned when a PostgreSQL repository cannot find an aggregate.
	ErrNotFound = sql.ErrNoRows
	// ErrUnsupported is returned by adapters whose persistence format is not yet
	// expressive enough to restore a rich domain object from storage.
	ErrUnsupported = errors.New("postgres adapter operation is unsupported")
)
