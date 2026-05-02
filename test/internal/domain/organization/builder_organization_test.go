package organization_test

import (
	"gitflic.ru/lms/internal/domain/organization"
	"gitflic.ru/lms/internal/domain/organization/inn"
	"gitflic.ru/lms/internal/domain/organization/orgname"
)

type organizationBuilder struct {
	name orgname.Name
	inn  inn.INN
}

func newOrganizationBuilder() *organizationBuilder {
	return &organizationBuilder{}
}

func (b *organizationBuilder) withName(name string) *organizationBuilder {
	b.name = makeName(name)
	return b
}

func (b *organizationBuilder) withInn(inn string, t inn.Type) *organizationBuilder {
	b.inn = makeINN(inn, t)
	return b
}

func (b *organizationBuilder) build() (*organization.Organization, error) {
	return organization.New(b.inn, b.name)
}

func makeName(s string) orgname.Name {
	name, _ := orgname.New(s)
	return name
}

func makeINN(s string, t inn.Type) inn.INN {
	inn, _ := inn.New(s, t)
	return inn
}
