package dob

import (
	"errors"
	"testing"
	"time"
)

func Test_validateAdultDateOfBirth(t *testing.T) {
	now := time.Date(2025, 6, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		date    time.Time
		asOf    time.Time
		wantErr bool
	}{
		{
			name:    "возраст меньше минимального",
			date:    time.Date(2020, 1, 12, 0, 0, 0, 0, time.UTC),
			asOf:    now,
			wantErr: true,
		},
		{
			name:    "возраст больше максимального",
			date:    time.Date(1800, 1, 12, 0, 0, 0, 0, time.UTC),
			asOf:    now,
			wantErr: true,
		},
		{
			name:    "возраст ровно минимальный порог",
			date:    now.AddDate(-MinAdultAge, 0, 0),
			asOf:    now,
			wantErr: false,
		},
		{
			name:    "возраст на один день меньше минимального",
			date:    now.AddDate(-MinAdultAge, 0, 1),
			asOf:    now,
			wantErr: true,
		},
		{
			name:    "возраст ровно максимальный порог",
			date:    now.AddDate(-MaxAdultAge, 0, 0),
			asOf:    now,
			wantErr: false,
		},
		{
			name:    "возраст на один год больше максимального",
			date:    now.AddDate(-(MaxAdultAge + 1), 0, 0),
			asOf:    now,
			wantErr: true,
		},
		{
			name:    "возраст в допустимом диапазоне",
			date:    time.Date(1995, 3, 10, 0, 0, 0, 0, time.UTC),
			asOf:    now,
			wantErr: false,
		},
		{
			name:    "дата рождения совпадает с asOf (возраст 0)",
			date:    now,
			asOf:    now,
			wantErr: true,
		},
		{
			name:    "дата рождения в будущем относительно asOf",
			date:    time.Date(2010, 1, 12, 0, 0, 0, 0, time.UTC),
			asOf:    time.Date(2002, 1, 12, 0, 0, 0, 0, time.UTC),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAdultDateOfBirth(tt.date, tt.asOf)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAdultDateOfBirth(%v, %v) error = %v, wantErr %v", tt.date, tt.asOf, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalid) {
				t.Errorf("validateAdultDateOfBirth(%v, %v) error %v не содержит ErrInvalid", tt.date, tt.asOf, err)
			}
		})
	}
}

func Test_validateMinAgeBoundary(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		wantErr bool
	}{
		{
			name:    "возраст на один год меньше минимального",
			age:     MinAdultAge - 1,
			wantErr: true,
		},
		{
			name:    "возраст равен минимальному порогу",
			age:     MinAdultAge,
			wantErr: false,
		},
		{
			name:    "возраст на один год больше минимального порога",
			age:     MinAdultAge + 1,
			wantErr: false,
		},
		{
			name:    "нулевой возраст",
			age:     0,
			wantErr: true,
		},
		{
			name:    "отрицательный возраст",
			age:     -1,
			wantErr: true,
		},
		{
			name:    "возраст в допустимом диапазоне меньше максимального",
			age:     MaxAdultAge - 1,
			wantErr: false,
		},
		{
			name: "возраст больше максимального - функция не возвращает ошибку",
			age:  MaxAdultAge + 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMinAgeBoundary(tt.age)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMinAgeBoundary(%d) error = %v, wantErr %v", tt.age, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalid) {
				t.Errorf("validateMinAgeBoundary(%d) error %v не содержит ErrInvalid", tt.age, err)
			}
		})
	}
}

func Test_validateMaxAgeBoundary(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		wantErr bool
	}{
		{
			name:    "возраст на один год больше максимального",
			age:     MaxAdultAge + 1,
			wantErr: true,
		},
		{
			name:    "возраст равен максимальному порогу",
			age:     MaxAdultAge,
			wantErr: false,
		},
		{
			name:    "возраст на один год меньше максимального порога",
			age:     MaxAdultAge - 1,
			wantErr: false,
		},
		{
			name:    "возраст значительно превышает максимальный",
			age:     MaxAdultAge + 100,
			wantErr: true,
		},
		{
			name:    "возраст в допустимом диапазоне больше минимального",
			age:     MinAdultAge + 1,
			wantErr: false,
		},
		{
			name:    "возраст меньше минимального - функция не возвращает ошибку",
			age:     MinAdultAge - 1,
			wantErr: false,
		},
		{
			name:    "возраст отрицательный - функция не возвращает ошибку",
			age:     -1,
			wantErr: false,
		},
		{
			name:    "нулевой возраст - функция не возвращает ошибку",
			age:     0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMaxAgeBoundary(tt.age)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateMaxAgeBoundary(%d) error = %v, wantErr %v", tt.age, err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalid) {
				t.Errorf("validateMaxAgeBoundary(%d) error %v не содержит ErrInvalid", tt.age, err)
			}
		})
	}
}
