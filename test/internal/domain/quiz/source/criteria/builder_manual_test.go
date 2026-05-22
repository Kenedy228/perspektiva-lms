//go:build legacy
// +build legacy

package criteria_test

import (
	"testing"

	criteria2 "gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type manualBuilder struct {
	questionIDs []uuid.UUID
}

func newManualBuilder() *manualBuilder {
	return &manualBuilder{
		questionIDs: []uuid.UUID{},
	}
}

func (b *manualBuilder) withQuestionID(id uuid.UUID) *manualBuilder {
	b.questionIDs = append(b.questionIDs, id)
	return b
}

func (b *manualBuilder) withMaxSizeQuestions() *manualBuilder {
	maxSize := int(1e5 + 1)
	questionIDs := make([]uuid.UUID, 0, maxSize)

	for range maxSize {
		questionIDs = append(questionIDs, uuid.New())
	}

	b.questionIDs = questionIDs
	return b
}

func (b *manualBuilder) build(t *testing.T, wantErr error) criteria2.Criteria {
	t.Helper()

	c, err := criteria2.NewManual(b.questionIDs)

	assert.ErrorIs(t, err, wantErr)

	return c
}

func (b *manualBuilder) buildNoTest() criteria2.Criteria {
	c, _ := criteria2.NewManual(b.questionIDs)
	return c
}
