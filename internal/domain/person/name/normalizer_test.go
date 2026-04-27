package name

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	var tests = map[string]struct {
		input string
		want  string
	}{
		"empty": {
			input: "",
			want:  "",
		},
		"spaces only": {
			input: "     ",
			want:  "",
		},
		"trim left and right": {
			input: "  Иван  ",
			want:  "Иван",
		},
		"collapse inner spaces": {
			input: "Иван   Петр",
			want:  "Иван Петр",
		},
		"collapse mixed whitespace": {
			input: "Иван\t\nПетр",
			want:  "Иван Петр",
		},
		"collapse with inner separator": {
			input: "Иван\t\t-\t\n Петров",
			want:  "Иван - Петров",
		},
		"already normalized": {
			input: "Иван Петр",
			want:  "Иван Петр",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			//Act
			got := normalize(tt.input)

			//Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
