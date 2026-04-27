package collection_test

import "github.com/google/uuid"

type fixture struct {
	id    uuid.UUID
	title string
}

func (f fixture) ID() uuid.UUID {
	return f.id
}
