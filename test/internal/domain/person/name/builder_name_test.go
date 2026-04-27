package name_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/person/name"
	"github.com/stretchr/testify/assert"
)

type nameBuilder struct {
	firstname  string
	lastname   string
	middlename string
}

func newNameBuilder() *nameBuilder {
	return &nameBuilder{
		firstname:  "",
		lastname:   "",
		middlename: "",
	}
}

func (b *nameBuilder) withFirstname(s string) *nameBuilder {
	b.firstname = s
	return b
}

func (b *nameBuilder) withLastname(s string) *nameBuilder {
	b.lastname = s
	return b
}

func (b *nameBuilder) withMiddlename(s string) *nameBuilder {
	b.middlename = s
	return b
}

func (b *nameBuilder) build(t *testing.T, wantErr error) name.Name {
	t.Helper()

	params := name.Params{
		Firstname:  b.firstname,
		Middlename: b.middlename,
		Lastname:   b.lastname,
	}

	n, err := name.New(params)

	assert.ErrorIs(t, err, wantErr)

	return n
}
