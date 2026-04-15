package base

import "github.com/google/uuid"

type Params struct {
	Text        string
	Description string
	Image       uuid.UUID
}
