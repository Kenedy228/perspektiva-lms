package dob

import (
	"testing"
	"time"
)

func Test_validateAdultDateOfBirth(t *testing.T) {
	type args struct {
		date time.Time
		asOf time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "возраст младше порогового значения",
			args: args{
				date: time.Date(2020, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "возраст старше порогового значения",
			args: args{
				date: time.Date(1800, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "валидный возраст",
			args: args{
				date: time.Date(2004, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			wantErr: false,
		},
		{
			name: "дата рождения позднее даты получения возраста",
			args: args{
				date: time.Date(2004, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Date(2002, 1, 12, 0, 0, 0, 0, time.UTC),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAdultDateOfBirth(tt.args.date, tt.args.asOf); (err != nil) != tt.wantErr {
				t.Errorf("validateAdultDateOfBirth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
