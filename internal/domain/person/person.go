package person

import (
	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Person struct {
	id      uuid.UUID
	name    name.Name
	profile *profile.Profile
}

func New(n name.Name) (*Person, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	return &Person{
		id:      id,
		name:    n,
		profile: nil,
	}, nil
}

func (p *Person) ID() uuid.UUID {
	return p.id
}

func (p *Person) Name() name.Name {
	return p.name
}

func (p *Person) Profile() (profile.Profile, bool) {
	if !p.HasProfile() {
		return profile.Profile{}, false
	}

	return *p.profile, true
}

func (p *Person) HasProfile() bool {
	if p.profile == nil {
		return false
	}

	return true
}

func (p *Person) AttachProfile(prof profile.Profile) {
	if p.HasProfile() {
		p.DetachProfile()
	}

	p.profile = &prof
}

func (p *Person) DetachProfile() {
	if !p.HasProfile() {
		return
	}

	p.profile = nil
}

func (p *Person) Rename(n name.Name) {
	p.name = n
}
