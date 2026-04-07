package person

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	type given struct {
		firstName    string
		lastName     string
		middleName   string
		organization uuid.UUID
	}

	type want struct {
		err error
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "empty firstName",
			given: given{
				firstName:    "",
				lastName:     "",
				middleName:   "",
				organization: uuid.Nil,
			},
			want: want{
				err: ErrEmptyFirstName,
			},
		},
		{
			name: "firstName with only spaces",
			given: given{
				firstName:    "    ",
				lastName:     "",
				middleName:   "",
				organization: uuid.Nil,
			},
			want: want{
				err: ErrEmptyFirstName,
			},
		},
		{
			name: "empty lastName",
			given: given{
				firstName:    "first",
				lastName:     "",
				middleName:   "",
				organization: uuid.Nil,
			},
			want: want{
				err: ErrEmptyLastName,
			},
		},
		{
			name: "lastName with only spaces",
			given: given{
				firstName:    "first",
				lastName:     "    ",
				middleName:   "",
				organization: uuid.Nil,
			},
			want: want{
				err: ErrEmptyLastName,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.given.firstName, tt.given.lastName, tt.given.middleName, tt.given.organization)

			if !errors.Is(err, tt.want.err) {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestRename(t *testing.T) {
	type given struct {
		firstName      string
		lastName       string
		middleName     string
		organizationID uuid.UUID
	}

	type when struct {
		firstName  string
		lastName   string
		middleName string
	}

	type want struct {
		firstName  string
		lastName   string
		middleName string
		err        error
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "empty first name",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:  "",
				lastName:   "newLast",
				middleName: "newMiddle",
			},
			want: want{
				firstName:  "first",
				lastName:   "last",
				middleName: "middle",
				err:        ErrEmptyFirstName,
			},
		},
		{
			name: "first name with only spaces",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:  "  ",
				lastName:   "newLast",
				middleName: "newMiddle",
			},
			want: want{
				firstName:  "first",
				lastName:   "last",
				middleName: "middle",
				err:        ErrEmptyFirstName,
			},
		},
		{
			name: "empty last name",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:  "newFirstName",
				lastName:   "",
				middleName: "newMiddle",
			},
			want: want{
				firstName:  "first",
				lastName:   "last",
				middleName: "middle",
				err:        ErrEmptyLastName,
			},
		},
		{
			name: "last name with only spaces",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:  "newFirstName",
				lastName:   "   ",
				middleName: "newMiddle",
			},
			want: want{
				firstName:  "first",
				lastName:   "last",
				middleName: "middle",
				err:        ErrEmptyLastName,
			},
		},
		{
			name: "valid new data",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:  "newFirstName",
				lastName:   "newLastName",
				middleName: "newMiddleName",
			},
			want: want{
				firstName:  "newFirstName",
				lastName:   "newLastName",
				middleName: "newMiddleName",
				err:        nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := New(tt.given.firstName, tt.given.lastName, tt.given.middleName, tt.given.organizationID)
			err := p.Rename(tt.when.firstName, tt.when.lastName, tt.when.middleName)

			if !errors.Is(err, tt.want.err) {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if p.firstName != tt.want.firstName {
				t.Fatalf("expected firstName %v, got %v", tt.want.firstName, p.firstName)
			}

			if p.lastName != tt.want.lastName {
				t.Fatalf("expected lastName %v, got %v", tt.want.lastName, p.lastName)
			}

			if p.middleName != tt.want.middleName {
				t.Fatalf("expected middleName %v, got %v", tt.want.middleName, p.middleName)
			}
		})
	}
}

func TestChangeOrganization(t *testing.T) {
	type given struct {
		firstName      string
		lastName       string
		middleName     string
		organizationID uuid.UUID
	}

	type when struct {
		organizationID uuid.UUID
	}

	type want struct {
		equality bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "same id",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				organizationID: uuid.Nil,
			},
			want: want{
				equality: true,
			},
		},
		{
			name: "new organization",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				organizationID: uuid.New(),
			},
			want: want{
				equality: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := New(tt.given.firstName, tt.given.lastName, tt.given.middleName, tt.given.organizationID)
			prevOrg := p.OrganizationID()
			p.ChangeOrganization(tt.when.organizationID)

			if eq := (p.organizationID == prevOrg); eq != tt.want.equality {
				t.Errorf("expected organization ID equal %v, got %v", tt.want.equality, eq)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type given struct {
		firstName      string
		lastName       string
		middleName     string
		organizationID uuid.UUID
	}

	type when struct {
		firstName      string
		lastName       string
		middleName     string
		organizationID uuid.UUID
	}

	type want struct {
		equal bool
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "different id",
			given: given{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			when: when{
				firstName:      "first",
				lastName:       "last",
				middleName:     "middle",
				organizationID: uuid.Nil,
			},
			want: want{
				equal: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenP, _ := New(tt.given.firstName, tt.given.lastName, tt.given.middleName, tt.given.organizationID)
			whenP, _ := New(tt.when.firstName, tt.when.lastName, tt.when.middleName, tt.when.organizationID)

			if eq := givenP.Equal(whenP); eq != tt.want.equal {
				t.Fatalf("expected equality %v, got %v", tt.want.equal, eq)
			}
		})
	}
}
