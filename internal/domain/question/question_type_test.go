package question

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTypeValid(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		isValid bool
	}{
		{
			name:    "invalid val",
			val:     "invalid",
			isValid: false,
		},
		{
			name:    "valid val",
			val:     "matching",
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tp := Type(tt.val)
			res := tp.IsValid()

			assert.Equal(t, res, tt.isValid)
			if res {
				assert.Equal(t, tp.String(), tt.val)
			}
		})
	}
}
