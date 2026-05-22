//go:build legacy
// +build legacy

package criteria_test

import (
	"testing"

	criteria2 "gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/stretchr/testify/assert"
)

type randomBuilder struct {
	questionCount int
}

func newRandomBuilder() *randomBuilder {
	return &randomBuilder{
		questionCount: 0,
	}
}

func (b *randomBuilder) withCount(count int) *randomBuilder {
	b.questionCount = count
	return b
}

func (b *randomBuilder) build(t *testing.T, wantErr error) criteria2.Criteria {
	t.Helper()

	c, err := criteria2.NewRandom(b.questionCount)

	assert.ErrorIs(t, err, wantErr)

	return c
}
