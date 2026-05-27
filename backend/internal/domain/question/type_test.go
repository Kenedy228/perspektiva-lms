package question

import "testing"

func TestType_DefaultInstruction(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "неизвестное значение возвращает пустую строку",
			t:    Type(""),
			want: "",
		},
		{
			name: "известное значение возвращает конкретную строку",
			t:    TypeSelectable,
			want: "выберите один или несколько правильных вариантов ответа",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.DefaultInstruction(); got != tt.want {
				t.Errorf("DefaultInstruction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want bool
	}{
		{
			name: "неизвестный тип возвращает false",
			t:    Type(""),
			want: false,
		},
		{
			name: "известный тип возвращает true",
			t:    TypeSelectable,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "возвращает известное значение как есть",
			t:    TypeSelectable,
			want: "selectable",
		},
		{
			name: "возвращает неизвестное значение как есть",
			t:    Type("undefined"),
			want: "undefined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_Title(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "для неизвестных значений возвращает пустую строку",
			t:    Type(""),
			want: "",
		},
		{
			name: "для известных значений возвращает непустую строку",
			t:    TypeSelectable,
			want: "выбор",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.Title(); got != tt.want {
				t.Errorf("Title() = %v, want %v", got, tt.want)
			}
		})
	}
}
