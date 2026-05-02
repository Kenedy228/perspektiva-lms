package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAgeAt(t *testing.T) {
	tc := map[string]struct {
		date time.Time
		asOf time.Time
		want int
	}{
		"новорождённый в тот же день": {
			date: time.Date(2000, 1, 10, 15, 4, 5, 0, time.UTC),
			asOf: time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			want: 0,
		},
		"день до дня рождения": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 9, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"в день рождения": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"день после дня рождения": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 5, 11, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"до дня рождения в том же году": {
			date: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 12, 30, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"после дня рождения в том же году": {
			date: time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"дата 29 февраля, тот же день в високосном году": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2020, 2, 29, 0, 0, 0, 0, time.UTC),
			want: 20,
		},
		"дата 29 февраля, невисокосный год, до псевдо-дня рождения": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 2, 27, 0, 0, 0, 0, time.UTC),
			want: 18,
		},
		"дата 29 февраля, невисокосный год, в псевдо-день рождения": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"дата 29 февраля, невисокосный год, после псевдо-дня рождения": {
			date: time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC),
			want: 19,
		},
		"дата оценки раньше года рождения": {
			date: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf: time.Date(1999, 5, 10, 0, 0, 0, 0, time.UTC),
			want: -1,
		},
	}

	for ttName, tt := range tc {
		t.Run(ttName, func(t *testing.T) {
			//Act
			got := ageAt(tt.date, tt.asOf)

			//Assert
			assert.Equal(t, tt.want, got)
		})
	}
}
