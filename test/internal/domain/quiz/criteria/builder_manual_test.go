package criteria_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
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

func (b *manualBuilder) build(t *testing.T, wantErr error) criteria.Criteria {
	t.Helper()

	c, err := criteria.NewManual(b.questionIDs)

	assert.ErrorIs(t, err, wantErr)

	return c
}

func (b *manualBuilder) buildNoTest() criteria.Criteria {
	c, _ := criteria.NewManual(b.questionIDs)
	return c
}
