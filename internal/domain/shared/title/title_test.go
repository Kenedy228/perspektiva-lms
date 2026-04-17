package title

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTitle(t *testing.T) {
	tests := []struct {
		name string
		val  string
		err  error
	}{
		{
			name: "empty title",
			val:  "",
			err:  ErrEmptyTitle,
		},
		{
			name: "whitespaces title",
			val:  " ",
			err:  ErrEmptyTitle,
		},
		{
			name: "valid title",
			val:  "valid",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewTitle(tt.val)

			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				assert.Equal(t, title.String(), tt.val)
			}
		})
	}
}
