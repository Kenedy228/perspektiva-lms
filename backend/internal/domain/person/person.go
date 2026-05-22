package person

import (
	"gitflic.ru/lms/backend/internal/domain/person/name"
	"gitflic.ru/lms/backend/internal/domain/person/profile"
	"gitflic.ru/lms/backend/internal/domain/shared/uid"
	"github.com/google/uuid"
)

// Person represents an individual participant in the LMS.
type Person struct {
	id      uuid.UUID
	name    name.Name
	profile *profile.Profile
}

// New creates a person with a generated ID.
func New(n name.Name) (*Person, error) {
	if err := validateName(n); err != nil {
		return nil, err
	}

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

// Restore recreates an existing person from persisted state.
func Restore(id uuid.UUID, n name.Name, prof *profile.Profile) (*Person, error) {
	if err := validateID(id); err != nil {
		return nil, err
	}

	if err := validateName(n); err != nil {
		return nil, err
	}

	if prof != nil {
		if err := validateProfile(*prof); err != nil {
			return nil, err
		}
	}

	return &Person{
		id:      id,
		name:    n,
		profile: prof,
	}, nil
}

// ID returns the person identifier.
func (p *Person) ID() uuid.UUID {
	return p.id
}

// Name returns the person's name.
func (p *Person) Name() name.Name {
	return p.name
}

// Profile returns the attached profile and whether it exists.
func (p *Person) Profile() (profile.Profile, bool) {
	if !p.HasProfile() {
		return profile.Profile{}, false
	}

	return *p.profile, true
}

// HasProfile reports whether a profile is attached.
func (p *Person) HasProfile() bool {
	if p.profile == nil {
		return false
	}

	return true
}

// AttachOrReplaceProfile attaches the person's initial profile.
func (p *Person) AttachOrReplaceProfile(prof profile.Profile) error {
	if err := validateProfile(prof); err != nil {
		return err
	}

	if p.HasProfile() {
		p.DetachProfile()
	}

	p.profile = &prof
	return nil
}

// DetachProfile removes the attached profile.
func (p *Person) DetachProfile() {
	if !p.HasProfile() {
		return
	}

	p.profile = nil
}

// ChangeName changes the person's name.
func (p *Person) ChangeName(n name.Name) error {
	if err := validateName(n); err != nil {
		return err
	}

	p.name = n
	return nil
}
