package snils_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/person/snils"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		input         string
		wantValue     string
		wantFormatted string
		wantErr       error
	}{
		"valid formatted input": {
			input:         "112-233-445 95",
			wantValue:     "11223344595",
			wantFormatted: "112-233-445 95",
			wantErr:       nil,
		},
		"valid plain input": {
			input:         "11223344595",
			wantValue:     "11223344595",
			wantFormatted: "112-233-445 95",
			wantErr:       nil,
		},
		"old snils below checksumFrom": {
			input:         "00100199800",
			wantValue:     "00100199800",
			wantFormatted: "001-001-998 00",
			wantErr:       nil,
		},
		"invalid checksum": {
			input:   "112-233-445 00",
			wantErr: snils.ErrInvalid,
		},
		"non digit in formatted": {
			input:   "112-233-445 9A",
			wantErr: snils.ErrInvalid,
		},
		"too short raw": {
			input:   "112-233-445 9",
			wantErr: snils.ErrInvalid,
		},
		"empty": {
			input:   "",
			wantErr: snils.ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			s := newSnilsBuilder().withValue(tt.input).
				build(t, tt.wantErr)

			//Assert
			assert.Equal(t, s.Value(), tt.wantValue)
			assert.Equal(t, s.Formatted(), tt.wantFormatted)
		})
	}
}
