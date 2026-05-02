package snils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tc := map[string]struct {
		in   string
		want string
	}{
		"ровно 11 цифр": {
			in:   "11223344595",
			want: "11223344595",
		},
		"формат с дефисами и пробелом": {
			in:   "112-233-445 95",
			want: "11223344595",
		},
		"лишние пробелы по краям": {
			in:   " 112-233-445 95 ",
			want: "11223344595",
		},
		"только разделители": {
			in:   "---   ---",
			want: "",
		},
		"пустая строка": {
			in:   "",
			want: "",
		},
	}

	for ttName, tt := range tc {
		t.Run(ttName, func(t *testing.T) {
			//Act
			normalized := normalize(tt.in)

			//Assert
			assert.Equal(t, normalized, tt.want)
		})
	}
}
