package profile

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"empty": {
			input: "",
			want:  "",
		},
		"only spaces": {
			input: "   \t   \n   ",
			want:  "",
		},
		"trim spaces": {
			input: "  инженер  ",
			want:  "инженер",
		},
		"collapse inner spaces": {
			input: "ведущий    инженер   программист",
			want:  "ведущий инженер программист",
		},
		"collapse tabs and newlines": {
			input: "высшее\tобразование\nмагистратура",
			want:  "высшее образование магистратура",
		},
		"already normalized": {
			input: "среднее профессиональное",
			want:  "среднее профессиональное",
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			got := normalize(tt.input)

			//Assert
			assert.Equal(t, got, tt.want)
		})
	}
}
