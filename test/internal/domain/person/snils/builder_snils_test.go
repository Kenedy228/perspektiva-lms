package snils_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/stretchr/testify/assert"
)

type snilsBuilder struct {
	value string
}

func newSnilsBuilder() *snilsBuilder {
	return &snilsBuilder{
		value: "",
	}
}

func (b *snilsBuilder) withValue(value string) *snilsBuilder {
	b.value = value
	return b
}

func (b *snilsBuilder) build(t *testing.T, wantErr error) snils.Snils {
	t.Helper()

	s, err := snils.New(b.value)

	assert.ErrorIs(t, err, wantErr)

	return s
}
