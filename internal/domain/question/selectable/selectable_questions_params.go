package selectable

import "github.com/google/uuid"

type Params struct {
	Text    string
	Image   uuid.UUID
	Options map[string]bool
}
