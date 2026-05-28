package limit

import (
	"errors"
	"testing"
	"time"
)

func TestNewAttempts(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{name: "ноль — бесконечные попытки", count: 0},
		{name: "одна попытка", count: 1},
		{name: "максимальное значение", count: maxAttemptsCount},
		{name: "отрицательное значение", count: -1, wantErr: true},
		{name: "превышение максимума", count: maxAttemptsCount + 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAttempts(tt.count)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got.Count() != tt.count {
				t.Fatalf("ожидалось Count()=%d, получено %d", tt.count, got.Count())
			}
		})
	}
}

func TestAttempts_IsInfinite(t *testing.T) {
	infinite, err := NewAttempts(0)
	if err != nil {
		t.Fatalf("create infinite attempts: %v", err)
	}
	if !infinite.IsInfinite() {
		t.Fatal("ожидалось IsInfinite()=true для 0")
	}

	finite, err := NewAttempts(5)
	if err != nil {
		t.Fatalf("create finite attempts: %v", err)
	}
	if finite.IsInfinite() {
		t.Fatal("ожидалось IsInfinite()=false для 5")
	}
}

func TestNewTime(t *testing.T) {
	tests := []struct {
		name    string
		seconds int
		wantErr bool
	}{
		{name: "ноль — без ограничения", seconds: 0},
		{name: "60 секунд", seconds: 60},
		{name: "максимальное значение", seconds: maxSecondsCount},
		{name: "отрицательное значение", seconds: -1, wantErr: true},
		{name: "превышение максимума", seconds: maxSecondsCount + 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewTime(tt.seconds)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got.Seconds() != tt.seconds {
				t.Fatalf("ожидалось Seconds()=%d, получено %d", tt.seconds, got.Seconds())
			}
		})
	}
}

func TestTime_IsInfinite(t *testing.T) {
	infinite, err := NewTime(0)
	if err != nil {
		t.Fatalf("create infinite time: %v", err)
	}
	if !infinite.IsInfinite() {
		t.Fatal("ожидалось IsInfinite()=true для 0")
	}

	finite, err := NewTime(600)
	if err != nil {
		t.Fatalf("create finite time: %v", err)
	}
	if finite.IsInfinite() {
		t.Fatal("ожидалось IsInfinite()=false для 600")
	}
}

func TestTime_TryDuration(t *testing.T) {
	infinite, _ := NewTime(0)
	if _, ok := infinite.TryDuration(); ok {
		t.Fatal("ожидалось ok=false для бесконечного ограничения")
	}

	finite, _ := NewTime(120)
	d, ok := finite.TryDuration()
	if !ok {
		t.Fatal("ожидалось ok=true для конечного ограничения")
	}
	if d != 120*time.Second {
		t.Fatalf("ожидалась длительность 120s, получено %v", d)
	}
}
