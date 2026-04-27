package organization_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/organization"
	"github.com/stretchr/testify/assert"
)

type organizationBuilder struct {
	name string
	inn  string
}

func newOrganizationBuilder() *organizationBuilder {
	return &organizationBuilder{
		name: "",
		inn:  "",
	}
}

func (b *organizationBuilder) withName(name string) *organizationBuilder {
	b.name = name
	return b
}

func (b *organizationBuilder) withInn(inn string) *organizationBuilder {
	b.inn = inn
	return b
}

func (b *organizationBuilder) build(t *testing.T, wantErr error) *organization.Organization {
	t.Helper()

	params := organization.Params{
		Name: b.name,
		INN:  b.inn,
	}

	org, err := organization.New(params)

	assert.ErrorIs(t, err, wantErr)

	return org
}
