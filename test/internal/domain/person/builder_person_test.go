package person_test

import (
	"gitflic.ru/lms/internal/domain/person"
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
)

type personBuilder struct {
	name    name.Name
	profile *profile.Profile
}

func newPersonBuilder() *personBuilder {
	return &personBuilder{}
}

func (b *personBuilder) withName() *personBuilder {
	b.name = nameFixture()
	return b
}

func (b *personBuilder) withProfile() *personBuilder {
	b.profile = profileFixture()
	return b
}

func (b *personBuilder) build() *person.Person {
	params := person.Params{
		Name:    b.name,
		Profile: b.profile,
	}

	person, _ := person.New(params)
	return person
}
