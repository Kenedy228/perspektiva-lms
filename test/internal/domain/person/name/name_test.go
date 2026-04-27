package name_test

import (
	"strings"
	"testing"

	"gitflic.ru/lms/internal/domain/person/name"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	var tests = map[string]struct {
		input      name.Params
		firstname  string
		lastname   string
		middlename string
		wantErr    error
	}{
		"valid without middlename": {
			input: name.Params{
				Firstname:  " Иван ",
				Lastname:   " Петров ",
				Middlename: " ",
			},
			firstname:  "Иван",
			lastname:   "Петров",
			middlename: "",
			wantErr:    nil,
		},
		"valid full fio": {
			input: name.Params{
				Firstname:  "Анна-Мария",
				Lastname:   "Иванова",
				Middlename: "Сергеевна",
			},
			firstname:  "Анна-Мария",
			lastname:   "Иванова",
			middlename: "Сергеевна",
			wantErr:    nil,
		},
		"valid collapses spaces": {
			input: name.Params{
				Firstname:  "Иван   Петр",
				Lastname:   "Сидоров",
				Middlename: "",
			},
			firstname:  "Иван Петр",
			lastname:   "Сидоров",
			middlename: "",
			wantErr:    nil,
		},
		"empty firstname invalid": {
			input: name.Params{
				Firstname:  "",
				Lastname:   "Петров",
				Middlename: "",
			},
			wantErr: name.ErrInvalid,
		},
		"empty lastname invalid": {
			input: name.Params{
				Firstname:  "Иван",
				Lastname:   "",
				Middlename: "",
			},
			wantErr: name.ErrInvalid,
		},
		"firstname too long invalid": {
			input: name.Params{
				Firstname:  strings.Repeat("А", 101),
				Lastname:   "Петров",
				Middlename: "",
			},
			wantErr: name.ErrInvalid,
		},
		"lastname digits invalid": {
			input: name.Params{
				Firstname:  "Иван",
				Lastname:   "Петр0в",
				Middlename: "",
			},
			wantErr: name.ErrInvalid,
		},
		"middlename invalid": {
			input: name.Params{
				Firstname:  "Иван",
				Lastname:   "Петров",
				Middlename: "Иваныч123",
			},
			wantErr: name.ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			n := newNameBuilder().withFirstname(tt.input.Firstname).
				withLastname(tt.input.Lastname).
				withMiddlename(tt.input.Middlename).
				build(t, tt.wantErr)

			//Assert
			assert.Equal(t, n.Firstname(), tt.firstname)
			assert.Equal(t, n.Lastname(), tt.lastname)
			assert.Equal(t, n.Middlename(), tt.middlename)
		})
	}
}

func TestFullname(t *testing.T) {
	var tests = map[string]struct {
		firstname  string
		lastname   string
		middlename string
		want       string
	}{
		"without middlename": {
			firstname:  "Иван",
			lastname:   "Петров",
			middlename: "",
			want:       "Петров Иван",
		},
		"with middlename": {
			firstname:  "Иван",
			lastname:   "Петров",
			middlename: "Иванович",
			want:       "Петров Иван Иванович",
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			n := newNameBuilder().withFirstname(tt.firstname).
				withLastname(tt.lastname).
				withMiddlename(tt.middlename).
				build(t, nil)

			//Act
			fullname := n.Fullname()

			//Assert
			assert.Equal(t, tt.want, fullname)
		})
	}
}

func TestWithInitials(t *testing.T) {
	var tests = map[string]struct {
		firstname  string
		lastname   string
		middlename string
		want       string
	}{
		"without middlename": {
			firstname:  "Иван",
			lastname:   "Петров",
			middlename: "",
			want:       "Петров И.",
		},
		"with middlename": {
			firstname:  "Иван",
			lastname:   "Петров",
			middlename: "Иванович",
			want:       "Петров И.И.",
		},
		"yo initials": {
			firstname:  "Ёж",
			lastname:   "Ёлкин",
			middlename: "Егорович",
			want:       "Ёлкин Ё.Е.",
		},
		"apostrophe lastname": {
			firstname:  "Иван",
			lastname:   "О'Коннор",
			middlename: "",
			want:       "О'Коннор И.",
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			n := newNameBuilder().withFirstname(tt.firstname).
				withLastname(tt.lastname).
				withMiddlename(tt.middlename).
				build(t, nil)

			//Act
			initials := n.WithInitials()

			//Assert
			assert.Equal(t, tt.want, initials)
		})
	}
}
