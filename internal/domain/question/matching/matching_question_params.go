package matching

import (
	"gitflic.ru/lms/internal/domain/content"
	"github.com/google/uuid"
)

type Params struct {
	Text       string
	Image      uuid.UUID
	Pairs      map[string]content.RichContent
	PairsCount int
}
