package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNormalize(t *testing.T) {
	tests := map[string]struct {
		input time.Time
		want  time.Time
	}{
		"zero": {
			input: time.Time{},
			want:  time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"non zero with time and zone": {
			input: time.Date(2000, 5, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)),
			want:  time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			normalized := normalize(tt.input)

			//Assert
			assert.Equal(t, normalized, tt.want)
		})
	}
}
