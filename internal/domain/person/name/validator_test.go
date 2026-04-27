package name

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredPart(t *testing.T) {
	var tests = map[string]struct {
		field   string
		input   string
		wantErr bool
	}{
		"simple valid": {
			field:   "имя",
			input:   "Иван",
			wantErr: false,
		},
		"two words valid": {
			field:   "имя",
			input:   "Иван Петр",
			wantErr: false,
		},
		"hyphen valid": {
			field:   "имя",
			input:   "Анна-Мария",
			wantErr: false,
		},
		"apostrophe valid": {
			field:   "фамилия",
			input:   "О'Коннор",
			wantErr: false,
		},
		"yo valid": {
			field:   "фамилия",
			input:   "Ёлкин",
			wantErr: false,
		},
		"max length valid": {
			field:   "имя",
			input:   strings.Repeat("А", 100),
			wantErr: false,
		},
		"empty invalid": {
			field:   "имя",
			input:   "",
			wantErr: true,
		},
		"spaces invalid": {
			field:   "имя",
			input:   "   ",
			wantErr: true,
		},
		"too long invalid": {
			field:   "имя",
			input:   strings.Repeat("А", 101),
			wantErr: true,
		},
		"digits only invalid": {
			field:   "имя",
			input:   "123",
			wantErr: true,
		},
		"contains digits invalid": {
			field:   "имя",
			input:   "Иван123",
			wantErr: true,
		},
		"latin invalid": {
			field:   "имя",
			input:   "Ivan",
			wantErr: true,
		},
		"underscore invalid": {
			field:   "имя",
			input:   "Иван_Петр",
			wantErr: true,
		},
		"comma invalid": {
			field:   "имя",
			input:   "Иван,Петр",
			wantErr: true,
		},
		"emoji invalid": {
			field:   "имя",
			input:   "Иван🙂",
			wantErr: true,
		},
		"leading hyphen invalid": {
			field:   "имя",
			input:   "-Иван",
			wantErr: true,
		},
		"trailing hyphen invalid": {
			field:   "имя",
			input:   "Иван-",
			wantErr: true,
		},
		"leading apostrophe invalid": {
			field:   "имя",
			input:   "'Иван",
			wantErr: true,
		},
		"trailing apostrophe invalid": {
			field:   "имя",
			input:   "Иван'",
			wantErr: true,
		},
		"double hyphen invalid": {
			field:   "имя",
			input:   "Иван--Петр",
			wantErr: true,
		},
		"double apostrophe invalid": {
			field:   "имя",
			input:   "Иван''Петр",
			wantErr: true,
		},
		"hyphen apostrophe invalid": {
			field:   "имя",
			input:   "Иван-'Петр",
			wantErr: true,
		},
		"apostrophe hyphen invalid": {
			field:   "имя",
			input:   "Иван'-Петр",
			wantErr: true,
		},
		"space after hyphen invalid": {
			field:   "имя",
			input:   "Иван- Петр",
			wantErr: true,
		},
		"space before hyphen invalid": {
			field:   "имя",
			input:   "Иван -Петр",
			wantErr: true,
		},
		"separator only hyphen invalid": {
			field:   "имя",
			input:   "-",
			wantErr: true,
		},
		"separator only apostrophe invalid": {
			field:   "имя",
			input:   "'",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//Act
			err := validateRequiredPart(tt.field, tt.input)

			//Assert
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalid)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOptionalPart(t *testing.T) {
	var tests = map[string]struct {
		field   string
		input   string
		wantErr bool
	}{
		"empty valid": {
			field:   "отчество",
			input:   "",
			wantErr: false,
		},
		"spaces valid after normalization": {
			field:   "отчество",
			input:   "",
			wantErr: false,
		},
		"simple valid": {
			field:   "отчество",
			input:   "Иванович",
			wantErr: false,
		},
		"max length valid": {
			field:   "отчество",
			input:   strings.Repeat("А", 100),
			wantErr: false,
		},
		"too long invalid": {
			field:   "отчество",
			input:   strings.Repeat("А", 101),
			wantErr: true,
		},
		"digits invalid": {
			field:   "отчество",
			input:   "Иваныч123",
			wantErr: true,
		},
		"leading hyphen invalid": {
			field:   "отчество",
			input:   "-Оглы",
			wantErr: true,
		},
		"trailing hyphen invalid": {
			field:   "отчество",
			input:   "Оглы-",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//Act
			err := validateOptionalPart(tt.field, tt.input)

			//Assert
			if tt.wantErr {
				assert.ErrorIs(t, err, ErrInvalid)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
