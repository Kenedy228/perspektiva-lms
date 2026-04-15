package short

import "github.com/google/uuid"

type Params struct {
	Text           string
	Image          uuid.UUID
	Answers        []string
	AllowDuplicate bool
}
