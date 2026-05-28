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
		{name: "корректное значение", input: "Банк вопросов по математике", want: "Банк вопросов по математике"},
		{name: "обрезает пробелы по краям", input: "  Банк  ", want: "Банк"},
		{name: "нормализует внутренние пробелы", input: "Банк  вопросов", want: "Банк вопросов"},
		{name: "пустая строка", input: "", wantErr: true},
		{name: "только пробелы", input: "   ", wantErr: true},
		{name: "превышение лимита символов", input: strings.Repeat("а", ValueCharsLimit+1), wantErr: true},
		{name: "ровно на границе лимита", input: strings.Repeat("а", ValueCharsLimit), want: strings.Repeat("а", ValueCharsLimit)},
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
				t.Fatal("ожидалось IsZero()=false для корректного заголовка")
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
