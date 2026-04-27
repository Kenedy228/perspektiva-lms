package limit_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/shared/limit"
	"github.com/stretchr/testify/assert"
)

type limitBuilder struct {
	value int
}

func newLimitBuilder() *limitBuilder {
	return &limitBuilder{
		value: 0,
	}
}

func (b *limitBuilder) withValue(value int) *limitBuilder {
	b.value = value
	return b
}

func (b *limitBuilder) build(t *testing.T, wantErr error) limit.Limit {
	limit, err := limit.New(b.value)
	assert.ErrorIs(t, err, wantErr)

	return limit
}
