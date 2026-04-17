package limit

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLimit(t *testing.T) {
	tests := []struct {
		name string
		val  int
		err  error
	}{
		{
			name: "negative val",
			val:  -1,
			err:  ErrInvalidLimit,
		},
		{
			name: "zero val",
			val:  0,
			err:  nil,
		},
		{
			name: "positive val",
			val:  1,
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Limit(tt.val)

			err := l.Validate()

			assert.ErrorIs(t, err, tt.err)
		})
	}
}

func TestIsInfifite(t *testing.T) {
	tests := []struct {
		name       string
		val        int
		isInfinite bool
	}{
		{
			name:       "infinite",
			val:        0,
			isInfinite: true,
		},
		{
			name:       "finite",
			val:        1,
			isInfinite: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Limit(tt.val)

			require.Nil(t, l.Validate())
			assert.Equal(t, l.IsInfinite(), tt.isInfinite)
		})
	}
}
