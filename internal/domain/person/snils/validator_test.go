package snils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLength(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"exact 11 digits": {
			value:   "11223344595",
			wantErr: nil,
		},
		"too short": {
			value:   "1122334459", // 10
			wantErr: ErrInvalid,
		},
		"too long": {
			value:   "112233445950", // 12
			wantErr: ErrInvalid,
		},
		"empty": {
			value:   "",
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateLength(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateContent(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"all digits ok": {
			value:   "11223344595",
			wantErr: nil,
		},
		"non digit symbol": {
			value:   "1122334459A",
			wantErr: ErrInvalid,
		},
		"three equal digits in a row": {
			value:   "11123344595", // три '1' подряд
			wantErr: ErrInvalid,
		},
		"two equal digits allowed": {
			value:   "11223344595", // максимум две одинаковых подряд
			wantErr: nil,
		},
		"empty": {
			value:   "",
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateContent(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateChecksum(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"below_checksum_from_valid_by_rule": {
			value:   "00000000001",
			wantErr: nil,
		},
		"just_below_checksum_from": {
			value:   "00100199700",
			wantErr: nil,
		},
		"equal_checksum_from": {
			value:   "00100199800",
			wantErr: nil,
		},
		"above_checksum_from_valid_checksum": {
			value:   "11223344595",
			wantErr: nil,
		},
		"above_checksum_from_invalid_checksum": {
			value:   "11223344500",
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateChecksum(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		value   string
		wantErr error
	}{
		"valid canonical": {
			value:   "11223344595",
			wantErr: nil,
		},
		"old canonical below checksumFrom": {
			value:   "00100199800",
			wantErr: nil,
		},
		"invalid checksum": {
			value:   "11223344500",
			wantErr: ErrInvalid,
		},
		"length_not_11": {
			value:   "1122334459",
			wantErr: ErrInvalid,
		},
		"non_digit": {
			value:   "1122334459A",
			wantErr: ErrInvalid,
		},
		"three_equal_digits": {
			value:   "11123344595",
			wantErr: ErrInvalid,
		},
		"empty": {
			value:   "",
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validate(tt.value)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
