package sequence

import (
	"gitflic.ru/lms/internal/domain/content"
	"github.com/google/uuid"
)

type Params struct {
	Text  string
	Image uuid.UUID
	Items []content.RichContent
}
