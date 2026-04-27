package profile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMaxLength(t *testing.T) {
	tests := map[string]struct {
		field   string
		value   string
		length  int
		wantErr error
	}{
		"shorter than limit": {
			field:   "должность",
			value:   "инженер",
			length:  1000,
			wantErr: nil,
		},
		"equal to limit": {
			field:   "должность",
			value:   strings.Repeat("а", 10),
			length:  10,
			wantErr: nil,
		},
		"greater than limit": {
			field:   "должность",
			value:   strings.Repeat("а", 11),
			length:  10,
			wantErr: ErrInvalid,
		},
		"empty allowed by this validator": {
			field:   "образование",
			value:   "",
			length:  10,
			wantErr: nil,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateMaxLength(tt.field, tt.value, tt.length)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateJobTitle(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"valid": {
			value:   "инженер",
			wantErr: nil,
		},
		"empty": {
			value:   "",
			wantErr: ErrInvalid,
		},
		"normalized empty should be validated by caller not here": {
			value:   "   ",
			wantErr: nil,
		},
		"equal max length": {
			value:   strings.Repeat("а", jobTitleCharsLimit),
			wantErr: nil,
		},
		"greater than max length": {
			value:   strings.Repeat("а", jobTitleCharsLimit+1),
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateJobTitle(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateEducation(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"valid": {
			value:   "высшее образование",
			wantErr: nil,
		},
		"empty": {
			value:   "",
			wantErr: ErrInvalid,
		},
		"equal max length": {
			value:   strings.Repeat("б", educationCharsLimit),
			wantErr: nil,
		},
		"greater than max length": {
			value:   strings.Repeat("б", educationCharsLimit+1),
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			err := validateEducation(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
