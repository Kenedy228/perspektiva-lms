package snils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"plain 11 digits": {
			in:   "11223344595",
			want: "11223344595",
		},
		"formatted with dashes and space": {
			in:   "112-233-445 95",
			want: "11223344595",
		},
		"extra spaces": {
			in:   " 112-233-445 95 ",
			want: "11223344595",
		},
		"only separators": {
			in:   "---   ---",
			want: "",
		},
		"empty": {
			in:   "",
			want: "",
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			normalized := normalize(tt.in)

			//Assert
			assert.Equal(t, normalized, tt.want)
		})
	}
}
