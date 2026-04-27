package person

import (
	"time"

	"gitflic.ru/lms/internal/domain/person/name"
	"gitflic.ru/lms/internal/domain/person/profile"
	"gitflic.ru/lms/internal/domain/shared/uid"
	"github.com/google/uuid"
)

type Person struct {
	id        uuid.UUID
	name      name.Name
	profile   *profile.Profile
	createdAt time.Time
	updatedAt time.Time
}

func New(params Params) (*Person, error) {
	id, err := uid.New()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	var cProfile *profile.Profile
	if params.Profile != nil {
		clone := params.Profile.Clone()
		cProfile = &clone
	}

	return &Person{
		id:        id,
		name:      params.Name,
		profile:   cProfile,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (p *Person) ID() uuid.UUID {
	return p.id
}

func (p *Person) Name() name.Name {
	return p.name
}

func (p *Person) Profile() (profile.Profile, bool) {
	if p.HasProfile() {
		return p.profile.Clone(), true
	}

	return profile.Profile{}, false
}

func (p *Person) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Person) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Person) HasProfile() bool {
	if p.profile == nil {
		return false
	}

	return true
}

func (p *Person) AttachProfile(profile *profile.Profile) {
	if !p.HasProfile() {
		p.DetachProfile()
	}

	cProfile := profile.Clone()
	p.profile = &cProfile
	p.updatedAt = time.Now()
}

func (p *Person) DetachProfile() {
	if !p.HasProfile() {
		return
	}

	p.profile = nil
	p.updatedAt = time.Now()
}

func (p *Person) Rename(name name.Name) {
	p.name = name
	p.updatedAt = time.Now()
}
