package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestValidateNotInFuture(t *testing.T) {
	tests := map[string]struct {
		date    time.Time
		asOf    time.Time
		wantErr error
	}{
		"same day ok": {
			date:    time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"past ok": {
			date:    time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"future invalid": {
			date:    time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: ErrInvalid,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateNotInFuture(tt.date, tt.asOf)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateAdultAge(t *testing.T) {
	tests := map[string]struct {
		date    time.Time
		asOf    time.Time
		wantErr error
	}{
		"exactly 18 on birthday": {
			date:    time.Date(2002, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 5, 10, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"just before 18 birthday": {
			date:    time.Date(2002, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 5, 9, 0, 0, 0, 0, time.UTC),
			wantErr: ErrInvalid,
		},
		"just after 18 birthday": {
			date:    time.Date(2002, 5, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 5, 11, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"much older than 18": {
			date:    time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"leap day under 18": {
			date:    time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2021, 2, 28, 0, 0, 0, 0, time.UTC), // 17
			wantErr: ErrInvalid,
		},
		"leap day exactly 18": {
			date:    time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2022, 2, 28, 0, 0, 0, 0, time.UTC), // 18 по твоему правилу
			wantErr: nil,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateAdultAge(tt.date, tt.asOf)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestValidateAdultDateOfBirth(t *testing.T) {
	tests := map[string]struct {
		date    time.Time
		asOf    time.Time
		wantErr error
	}{
		"valid adult today": {
			date:    time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
		"underage today": {
			date:    time.Date(2005, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: ErrInvalid,
		},
		"future dob invalid": {
			date:    time.Date(2030, 1, 10, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
			wantErr: ErrInvalid,
		},
		"leap dob adult non leap asOf": {
			date:    time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2019, 2, 28, 0, 0, 0, 0, time.UTC),
			wantErr: nil,
		},
	}

	for ttName, tt := range tests {
		t.Run(ttName, func(t *testing.T) {
			//Act
			err := validateAdultDateOfBirth(tt.date, tt.asOf)

			//Assert
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
