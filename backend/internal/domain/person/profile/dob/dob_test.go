package dob

import (
	"reflect"
	"testing"
	"time"
)

func TestDateOfBirth_Date(t *testing.T) {
	type fields struct {
		date time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{
			name: "возвращает дату как есть",
			fields: fields{
				date: time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(2020, 2, 22, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DateOfBirth{
				date: tt.fields.date,
			}
			if got := db.Date(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Date() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDateOfBirth_IsZero(t *testing.T) {
	type fields struct {
		date time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустое значение",
			fields: fields{
				date: time.Time{},
			},
			want: true,
		},
		{
			name: "непустое значение",
			fields: fields{
				date: time.Now(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := DateOfBirth{
				date: tt.fields.date,
			}
			if got := db.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
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
			name: "возраст младше порогового значения",
			args: args{
				date: time.Date(2020, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			want: DateOfBirth{
				date: time.Time{},
			},
			wantErr: true,
		},
		{
			name: "возраст старше порогового значения",
			args: args{
				date: time.Date(1800, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			want: DateOfBirth{
				date: time.Time{},
			},
			wantErr: true,
		},
		{
			name: "валидный возраст",
			args: args{
				date: time.Date(2004, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Now(),
			},
			want: DateOfBirth{
				date: time.Date(2004, 1, 12, 0, 0, 0, 0, time.UTC),
			},
			wantErr: false,
		},
		{
			name: "дата рождения позднее даты получения возраста",
			args: args{
				date: time.Date(2004, 1, 12, 0, 0, 0, 0, time.UTC),
				asOf: time.Date(2002, 1, 12, 0, 0, 0, 0, time.UTC),
			},
			want: DateOfBirth{
				date: time.Time{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.date, tt.args.asOf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
