package dob_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"github.com/stretchr/testify/assert"
)

type dobBuilder struct {
	date time.Time
	asOf time.Time
}

func newDobBuilder() *dobBuilder {
	return &dobBuilder{
		date: time.Now(),
		asOf: time.Now(),
	}
}

func (b *dobBuilder) withDate(date time.Time) *dobBuilder {
	b.date = date
	return b
}

func (b *dobBuilder) withAsOf(asOf time.Time) *dobBuilder {
	b.asOf = asOf
	return b
}

func (b *dobBuilder) build(t *testing.T, wantErr error) dob.DateOfBirth {
	t.Helper()

	db, err := dob.New(b.date, b.asOf)

	assert.ErrorIs(t, err, wantErr)

	return db
}
