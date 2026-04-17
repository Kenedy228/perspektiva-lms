package s3validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name string
		key  string
		err  error
	}{
		{
			name: "empty key",
			key:  "",
			err:  ErrEmptyS3Key,
		},
		{
			name: "key contains space",
			key:  "k ey",
			err:  ErrInvalidS3Key,
		},
		{
			name: "key contains control",
			key:  "k\tey",
			err:  ErrInvalidS3Key,
		},
		{
			name: "valid key",
			key:  "images/image.png",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateS3Key(tt.key)

			assert.ErrorIs(t, tt.err, err)
		})
	}
}
