package question

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQText(t *testing.T) {
	tests := []struct {
		name string
		val  string
		err  error
	}{
		{
			name: "empty text",
			val:  "",
			err:  ErrEmptyText,
		},
		{
			name: "whitespaces text",
			val:  " ",
			err:  ErrEmptyText,
		},
		{
			name: "valid text",
			val:  "valid",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewQText(tt.val)

			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				assert.Equal(t, title.String(), tt.val)
			}
		})
	}
}
