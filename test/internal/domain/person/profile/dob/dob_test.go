package dob_test

import (
	"testing"
	"time"

	"gitflic.ru/lms/internal/domain/person/profile/dob"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("успех", func(t *testing.T) {
		tc := []struct {
			name string
			date time.Time
			asOf time.Time
		}{
			{
				name: "корректный возраст взрослого человека (18+)",
				date: time.Date(2000, 1, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)),
				asOf: time.Date(2020, 1, 10, 10, 0, 0, 0, time.UTC),
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				db, err := dob.New(tt.date, tt.asOf)
				expectedDate := time.Date(tt.date.Year(), tt.date.Month(), tt.date.Day(), 0, 0, 0, 0, time.UTC)

				// Assert
				assert.NoError(t, err)
				assert.Equal(t, expectedDate, db.Date())
			})
		}
	})

	t.Run("ошибка", func(t *testing.T) {
		tc := []struct {
			name    string
			date    time.Time
			asOf    time.Time
			wantErr error
		}{
			{
				name:    "возраст молодого человека (до 18)",
				date:    time.Date(2005, 1, 10, 0, 0, 0, 0, time.UTC),
				asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
				wantErr: dob.ErrInvalid,
			},
			{
				name:    "родился в будущем",
				date:    time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC),
				asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
				wantErr: dob.ErrInvalid,
			},
		}

		for _, tt := range tc {
			t.Run(tt.name, func(t *testing.T) {
				// Arrange
				_, err := dob.New(tt.date, tt.asOf)

				// Assert
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})
}
