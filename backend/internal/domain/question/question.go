package question

import (
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
	"github.com/google/uuid"
)

type Question interface {
	ID() uuid.UUID
	Title() title.Title
	Attachment() (attachment.Attachment, bool)
	Instruction() string
	Type() Type
	Clone() Question
	ChangeTitle(title.Title)
	ChangeAttachment(attachment.Attachment)
	RemoveAttachment()
	HasAttachment() bool
}
