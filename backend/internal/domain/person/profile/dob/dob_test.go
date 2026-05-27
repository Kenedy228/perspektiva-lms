package dob

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	asOf := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)

	type args struct {
		date time.Time
		asOf time.Time
	}

	tests := []struct {
		name    string
		args    args
		want    DateOfBirth
		wantErr bool
	}{
		{
			name: "ровно 18 лет",
			args: args{
				date: asOf.AddDate(-MinAdultAge, 0, 0),
				asOf: asOf,
			},
			want:    DateOfBirth{date: normalize(asOf.AddDate(-MinAdultAge, 0, 0))},
			wantErr: false,
		},
		{
			name: "17 лет и 1 день — на день младше нижней границы возраста",
			args: args{
				date: asOf.AddDate(-MinAdultAge, 0, 1),
				asOf: asOf,
			},
			want:    DateOfBirth{},
			wantErr: true,
		},
		{
			name: "18 лет и 1 день - на день старше нижней границы возраста",
			args: args{
				date: asOf.AddDate(-MinAdultAge, 0, -1),
				asOf: asOf,
			},
			want:    DateOfBirth{date: normalize(asOf.AddDate(-MinAdultAge, 0, -1))},
			wantErr: false,
		},
		{
			name: "достигнута верхняя граница возраста",
			args: args{
				date: asOf.AddDate(-MaxAdultAge, 0, 0),
				asOf: asOf,
			},
			want:    DateOfBirth{date: normalize(asOf.AddDate(-MaxAdultAge, 0, 0))},
			wantErr: false,
		},
		{
			name: "выше верхней границы возраста",
			args: args{
				date: asOf.AddDate(-(MaxAdultAge + 1), 0, 0),
				asOf: asOf,
			},
			want:    DateOfBirth{},
			wantErr: true,
		},
		{
			name: "ниже верхней границы возраста",
			args: args{
				date: asOf.AddDate(-(MaxAdultAge - 1), 0, 0),
				asOf: asOf,
			},
			want:    DateOfBirth{date: normalize(asOf.AddDate(-(MaxAdultAge - 1), 0, 0))},
			wantErr: false,
		},
		{
			name: "дата рождения в будущем относительно asOf",
			args: args{
				date: asOf.AddDate(0, 0, 1),
				asOf: asOf,
			},
			want:    DateOfBirth{},
			wantErr: true,
		},
		{
			name: "дата совпадает с asOf (возраст 0)",
			args: args{
				date: asOf,
				asOf: asOf,
			},
			want:    DateOfBirth{},
			wantErr: true,
		},
		{
			name: "нулевая дата рождения",
			args: args{
				date: time.Time{},
				asOf: asOf,
			},
			want:    DateOfBirth{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.date, tt.args.asOf)

			if tt.wantErr {
				assert.Error(t, err)
				assert.ErrorIs(t, err, ErrInvalid)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDateOfBirth_Date(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
	}{
		{
			name: "возвращает дату как есть",
			date: time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DateOfBirth{date: tt.date}
			assert.True(t, db.Date().Equal(tt.date),
				"Date() = %v, want %v", db.Date(), tt.date)
		})
	}
}

func TestDateOfBirth_IsZero(t *testing.T) {
	tests := []struct {
		name string
		date time.Time
		want bool
	}{
		{
			name: "пустое значение",
			date: time.Time{},
			want: true,
		},
		{
			name: "непустое значение",
			date: time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DateOfBirth{date: tt.date}
			assert.Equal(t, tt.want, db.IsZero())
		})
	}
}
