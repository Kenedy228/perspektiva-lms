package option

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewContentOption(t *testing.T) {
	tests := []struct {
		name  string
		cType ContentType
		val   string
		err   error
	}{
		{
			name:  "invalid cType",
			cType: ContentType("invalid"),
			val:   "text",
			err:   ErrInvalidContentType,
		},
		{
			name:  "invalid val, cType text",
			cType: ContentTypeText,
			val:   "",
			err:   ErrInvalidValue,
		},
		{
			name:  "invalid whitespaces val, cType text",
			cType: ContentTypeText,
			val:   " ",
			err:   ErrInvalidValue,
		},
		{
			name:  "invalid s3 key val, cType image",
			cType: ContentTypeImage,
			val:   "s3\t",
			err:   ErrInvalidValue,
		},
		{
			name:  "valid text value",
			cType: ContentTypeText,
			val:   "text",
			err:   nil,
		},
		{
			name:  "valid image value",
			cType: ContentTypeImage,
			val:   "images/s3",
			err:   nil,
		},
		{
			name:  "valid audio value",
			cType: ContentTypeAudio,
			val:   "audios/s3",
			err:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := NewContentOption(tt.cType, tt.val)

			assert.ErrorIs(t, err, tt.err)

			if err == nil {
				assert.Equal(t, opt.ContentType(), tt.cType)
				assert.Equal(t, opt.Value(), tt.val)
			}
		})
	}
}

func TestEqualContentOption(t *testing.T) {
	base, err := NewContentOption(ContentTypeText, "text")
	require.Nil(t, err)

	tests := []struct {
		name  string
		cType ContentType
		val   string
		equal bool
	}{
		{
			name:  "same",
			cType: ContentTypeText,
			val:   "text",
			equal: true,
		},
		{
			name:  "different content types",
			cType: ContentTypeImage,
			val:   "text",
			equal: false,
		},
		{
			name:  "different values",
			cType: ContentTypeText,
			val:   "another text",
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opt, err := NewContentOption(tt.cType, tt.val)
			require.Nil(t, err)

			if tt.equal {
				assert.Equal(t, opt, base)
			} else {
				assert.NotEqual(t, opt, base)
			}
		})
	}
}
