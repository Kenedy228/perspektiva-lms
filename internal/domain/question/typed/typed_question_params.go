package typed

import "github.com/google/uuid"

type Params struct {
	Text              string
	Image             uuid.UUID
	PlaceholdersCount int
	Blanks            map[string][]string
}
