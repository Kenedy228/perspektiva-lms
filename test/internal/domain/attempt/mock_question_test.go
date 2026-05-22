//go:build legacy
// +build legacy

package attempt_test

import (
	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/attachment"
	"gitflic.ru/lms/backend/internal/domain/question/title"
	"github.com/google/uuid"
)

type mockQuestion struct {
	id uuid.UUID
}

func (m mockQuestion) ID() uuid.UUID {
	return m.id
}

func (m mockQuestion) Clone() question2.Question {
	return m
}

func (m mockQuestion) Title() title.Title {
	panic("")
}

func (m mockQuestion) Attachment() (attachment.Attachment, bool) {
	panic("")
}

func (m mockQuestion) Instruction() string {
	panic("")
}

func (m mockQuestion) Type() question2.Type {
	panic("")
}

func (m mockQuestion) ChangeTitle(t title.Title) {
	panic("")
}

func (m mockQuestion) ChangeAttachment(a attachment.Attachment) {
	panic("")
}

func (m mockQuestion) RemoveAttachment() {
	panic("")
}

func (m mockQuestion) HasAttachment() bool {
	panic("")
}
