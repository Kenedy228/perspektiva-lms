package account

import (
	"testing"

	"github.com/google/uuid"
)

func TestNew(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
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
			name: "should create",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "empty login",
			given: given{
				login:    "",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				err: ErrEmptyLogin,
			},
		},
		{
			name: "empty hash",
			given: given{
				login:    "login",
				hash:     "",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				err: ErrEmptyPasswordHash,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)

			if err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}
		})
	}
}

func TestBlock(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
	}

	type want struct {
		blocked bool
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "blocked",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				blocked: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, _ := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)

			acc.Block()

			if acc.IsBlocked() != tt.want.blocked {
				t.Fatalf("expected blocked %v, got %v", tt.want.blocked, acc.IsBlocked())
			}
		})
	}
}

func TestUnblock(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
	}

	type want struct {
		blocked bool
	}

	tests := []struct {
		name string
		given
		want
	}{
		{
			name: "not blocked",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				blocked: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, _ := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)

			acc.Unblock()

			if acc.IsBlocked() != tt.want.blocked {
				t.Fatalf("expected blocked %v, got %v", tt.want.blocked, acc.IsBlocked())
			}
		})
	}
}

func TestEqual(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
	}

	type when struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
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
			name: "not equal",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			when: when{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			want: want{
				equal: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			givenAcc, _ := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)
			whenAcc, _ := New(tt.when.login, tt.when.hash, tt.when.roleID, tt.when.personID)

			if eq := givenAcc.Equal(whenAcc); eq != tt.want.equal {
				t.Fatalf("expected acc equal %v, got %v", tt.want.equal, eq)
			}
		})
	}
}

func TestChangeLogin(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
	}

	type when struct {
		login string
	}

	type want struct {
		err   error
		login string
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "empty login",
			given: given{
				login:    "old login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			when: when{
				login: "",
			},
			want: want{
				err:   ErrEmptyLogin,
				login: "old login",
			},
		},
		{
			name: "new login",
			given: given{
				login:    "old login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			when: when{
				login: "new login",
			},
			want: want{
				err:   nil,
				login: "new login",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, _ := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)

			err := acc.ChangeLogin(tt.when.login)

			if err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if acc.Login() != tt.want.login {
				t.Fatalf("expected login %v, got %v", tt.want.login, acc.Login())
			}
		})
	}
}

func TestChangePassword(t *testing.T) {
	type given struct {
		login    string
		hash     string
		roleID   uuid.UUID
		personID uuid.UUID
	}

	type when struct {
		hash string
	}

	type want struct {
		err  error
		hash string
	}

	tests := []struct {
		name string
		given
		when
		want
	}{
		{
			name: "empty hash",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			when: when{
				hash: "",
			},
			want: want{
				err:  ErrEmptyPasswordHash,
				hash: "hash",
			},
		},
		{
			name: "new hash",
			given: given{
				login:    "login",
				hash:     "hash",
				roleID:   uuid.New(),
				personID: uuid.New(),
			},
			when: when{
				hash: "new hash",
			},
			want: want{
				err:  nil,
				hash: "new hash",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, _ := New(tt.given.login, tt.given.hash, tt.given.roleID, tt.given.personID)

			err := acc.ChangePassword(tt.when.hash)

			if err != tt.want.err {
				t.Fatalf("expected err %v, got %v", tt.want.err, err)
			}

			if acc.passwordHash != tt.want.hash {
				t.Fatalf("expected hash %v, got %v", tt.want.hash, acc.passwordHash)
			}
		})
	}
}
