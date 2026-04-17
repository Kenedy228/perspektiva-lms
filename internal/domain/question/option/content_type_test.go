package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidContentType(t *testing.T) {
	tests := []struct {
		name    string
		val     string
		isValid bool
	}{
		{
			name:    "invalid",
			val:     "invalid",
			isValid: false,
		},
		{
			name:    "valid",
			val:     "image",
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ct := ContentType(tt.val)

			valid := ct.IsValid()

			assert.Equal(t, tt.isValid, valid)

			if valid {
				assert.Equal(t, tt.val, ct.String())
			}
		})
	}
}

func TestEqualContentType(t *testing.T) {
	tests := []struct {
		name  string
		cur   ContentType
		other ContentType
		equal bool
	}{
		{
			name:  "same",
			cur:   ContentTypeText,
			other: ContentTypeText,
			equal: true,
		},
		{
			name:  "different",
			cur:   ContentTypeText,
			other: ContentTypeImage,
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.equal {
				assert.Equal(t, tt.cur, tt.other)
			} else {
				assert.NotEqual(t, tt.cur, tt.other)
			}
		})
	}
}
