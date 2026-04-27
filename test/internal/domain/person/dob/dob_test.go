package dob_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person/dob"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := map[string]struct {
		date    time.Time
		asOf    time.Time
		want    time.Time
		wantErr error
	}{
		"valid adult": {
			date:    time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)),
			asOf:    time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC),
			want:    time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"underage error": {
			date:    time.Date(2005, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: dob.ErrInvalid,
		},
		"future dob error": {
			date:    time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: dob.ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Arrange
			db := newDobBuilder().withDate(tt.date).
				withAsOf(tt.asOf).
				build(t, tt.wantErr)

			//Assert
			assert.Equal(t, db.Date(), tt.want)
		})
	}
}
