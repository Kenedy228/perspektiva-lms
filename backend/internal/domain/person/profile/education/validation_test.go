package education

import (
	"strings"
	"testing"
)

func Test_validateValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустое значение",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "один ASCII символ",
			args: args{
				value: "a",
			},
			wantErr: false,
		},
		{
			name: "короткая строка",
			args: args{
				value: strings.Repeat("a", 20),
			},
			wantErr: false,
		},
		{
			name: "ровно лимит ASCII символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "на один символ больше лимита",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "значительно превышает лимит",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit*2),
			},
			wantErr: true,
		},
		{
			name: "пробельные символы в пределах лимита",
			args: args{
				value: strings.Repeat(" ", 20),
			},
			wantErr: false,
		},
		{
			name: "пробельные символы ровно лимит",
			args: args{
				value: strings.Repeat(" ", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "пробельные символы за лимитом",
			args: args{
				value: strings.Repeat(" ", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "unicode строка в пределах лимита",
			args: args{
				value: strings.Repeat("Я", 100),
			},
			wantErr: false,
		},
		{
			name: "unicode строка ровно лимит",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "unicode строка за лимитом",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "эмодзи в пределах лимита",
			args: args{
				value: strings.Repeat("🎓", 50),
			},
			wantErr: false,
		},
		{
			name: "эмодзи за лимитом",
			args: args{
				value: strings.Repeat("🎓", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "смешанные ASCII и unicode",
			args: args{
				value: strings.Repeat("a", 250) + strings.Repeat("Я", 250),
			},
			wantErr: false,
		},
		{
			name: "смешанные символы за лимитом",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit/2) + strings.Repeat("Я", ValueCharsLimit/2+1),
			},
			wantErr: true,
		},
		{
			name: "только табуляции",
			args: args{
				value: "\t\t\t\t\t",
			},
			wantErr: false,
		},
		{
			name: "только переводы строк",
			args: args{
				value: "\n\n\n\n\n",
			},
			wantErr: false,
		},
		{
			name: "смешанные пробельные символы",
			args: args{
				value: "\t \n  \r",
			},
			wantErr: false,
		},
		{
			name: "нулевые символы",
			args: args{
				value: "\x00\x00\x00",
			},
			wantErr: false,
		},
		{
			name: "символы нулевой ширины",
			args: args{
				value: "\u200B\u200B\u200B",
			},
			wantErr: false,
		},
		{
			name: "строка с эмодзи и текстом",
			args: args{
				value: "Высшее образование 🎓 Московский государственный университет",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateRequiredValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустая строка",
			args: args{
				value: "",
			},
			wantErr: true,
		},
		{
			name: "один символ",
			args: args{
				value: "a",
			},
			wantErr: false,
		},
		{
			name: "строка из нескольких символов",
			args: args{
				value: "hello world",
			},
			wantErr: false,
		},
		{
			name: "только пробелы",
			args: args{
				value: "     ",
			},
			wantErr: false,
		},
		{
			name: "только табуляции",
			args: args{
				value: "\t\t\t",
			},
			wantErr: false,
		},
		{
			name: "только переводы строк",
			args: args{
				value: "\n\n\n",
			},
			wantErr: false,
		},
		{
			name: "смешанные пробельные символы",
			args: args{
				value: "\t \n  \r",
			},
			wantErr: false,
		},
		{
			name: "unicode строка",
			args: args{
				value: "Привет, мир!",
			},
			wantErr: false,
		},
		{
			name: "строка с эмодзи",
			args: args{
				value: "образование 🎓",
			},
			wantErr: false,
		},
		{
			name: "нулевой символ",
			args: args{
				value: "\x00",
			},
			wantErr: false,
		},
		{
			name: "символы нулевой ширины",
			args: args{
				value: "\u200B",
			},
			wantErr: false,
		},
		{
			name: "длинная строка",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "строка превышает лимит символов - функция не возвращает ошибку",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequiredValue(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateValueCharsLimit(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустая строка - функция не возвращает ошибку",
			args: args{
				value: "",
			},
			wantErr: false,
		},
		{
			name: "один ASCII символ",
			args: args{
				value: "a",
			},
			wantErr: false,
		},
		{
			name: "ровно лимит ASCII символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "на один символ больше лимита",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "значительно превышает лимит",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit*2),
			},
			wantErr: true,
		},
		{
			name: "unicode символы в пределах лимита",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "unicode символы за лимитом",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "эмодзи в пределах лимита",
			args: args{
				value: strings.Repeat("🎓", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "эмодзи за лимитом",
			args: args{
				value: strings.Repeat("🎓", ValueCharsLimit+1),
			},
			wantErr: true,
		},
		{
			name: "смешанные ASCII и unicode в пределах лимита",
			args: args{
				value: strings.Repeat("a", 250) + strings.Repeat("Я", 250),
			},
			wantErr: false,
		},
		{
			name: "пробельные символы",
			args: args{
				value: strings.Repeat(" ", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "нулевые символы в пределах лимита",
			args: args{
				value: strings.Repeat("\x00", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "символы нулевой ширины",
			args: args{
				value: strings.Repeat("\u200B", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "нулевые символы за лимитом",
			args: args{
				value: strings.Repeat("\x00", ValueCharsLimit+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateValueCharsLimit(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("validateCharsLimit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
