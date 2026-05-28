package title

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "корректное значение", input: "Введение в Go", want: "Введение в Go"},
		{name: "пустая строка", input: "", wantErr: true},
		{name: "только пробелы", input: "   ", wantErr: true},
		{name: "ровно на границе лимита", input: strings.Repeat("а", valueCharsLimit), want: strings.Repeat("а", valueCharsLimit)},
		{name: "превышение лимита символов", input: strings.Repeat("а", valueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.input)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got.Value() != tt.want {
				t.Fatalf("ожидалось %q, получено %q", tt.want, got.Value())
			}
			if got.IsZero() {
				t.Fatal("ожидалось IsZero()=false для корректного названия")
			}
		})
	}
}

func TestTitle_IsZero(t *testing.T) {
	var zero Title
	if !zero.IsZero() {
		t.Fatal("ожидалось IsZero()=true для нулевого значения")
	}
}
