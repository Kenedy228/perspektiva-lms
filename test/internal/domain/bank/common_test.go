//go:build legacy
// +build legacy

package bank_test

import "github.com/google/uuid"

func makeIDs(size int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, size)

	for range size {
		ids = append(ids, uuid.New())
	}

	return ids
}
