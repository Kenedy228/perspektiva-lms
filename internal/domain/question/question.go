package question

import (
	"gitflic.ru/lms/internal/domain/question/attachment"
	"gitflic.ru/lms/internal/domain/question/title"
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
