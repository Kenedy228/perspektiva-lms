package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgeAt(t *testing.T) {
	tests := map[string]struct {
		date time.Time
		asOf time.Time
		want int
	}{
		"same day newborn": {
			date: time.Date(2000, 1, 10, 15, 4, 5, 0, time.UTC),
			asOf: time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		"day before birthday": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 9, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"on birthday": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"day after birthday": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 11, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"before birthday same year": {
			date: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"after birthday same year": {
			date: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"leap day to leap year same day": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"leap day to non-leap before pseudo birthday": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 2, 27, 0, 0, 0, 0, time.UTC),
			want: 18,
		},
		"leap day to non-leap on pseudo birthday": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"leap day to non-leap after pseudo birthday": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"asOf earlier year": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(1999, 5, 10, 0, 0, 0, 0, time.UTC),
			want: -1,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			got := ageAt(tt.date, tt.asOf)

			//Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
