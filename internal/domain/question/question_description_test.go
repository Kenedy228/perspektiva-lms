package question

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQDescription(t *testing.T) {
	tests := []struct {
		name string
		val  string
		err  error
	}{
		{
			name: "empty description",
			val:  "",
			err:  ErrEmptyDescription,
		},
		{
			name: "whitespaces description",
			val:  " ",
			err:  ErrEmptyDescription,
		},
		{
			name: "valid description",
			val:  "valid",
			err:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewQDescription(tt.val)

			assert.ErrorIs(t, err, tt.err)
			if err == nil {
				assert.Equal(t, title.String(), tt.val)
			}
		})
	}
}
