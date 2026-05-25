package education

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Education
		wantErr bool
	}{
		{
			name: "пустая строка",
			args: args{
				value: "",
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "только пробелы",
			args: args{
				value: strings.Repeat(" ", 20),
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "только табуляции",
			args: args{
				value: "\t\t\t",
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "только переводы строк",
			args: args{
				value: "\n\n\n",
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "смешанные пробельные символы",
			args: args{
				value: " \t \n  \r",
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "один ASCII символ",
			args: args{
				value: "a",
			},
			want: Education{
				value: "a",
			},
			wantErr: false,
		},
		{
			name: "короткая строка без пробелов",
			args: args{
				value: strings.Repeat("a", 20),
			},
			want: Education{
				value: strings.Repeat("a", 20),
			},
			wantErr: false,
		},
		{
			name: "ровно лимит ASCII символов",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: Education{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "на один символ больше лимита",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit+1),
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "значительно превышает лимит",
			args: args{
				value: strings.Repeat("a", ValueCharsLimit*2),
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "ведущие и замыкающие пробелы",
			args: args{
				value: "  высшее образование  ",
			},
			want: Education{
				value: "высшее образование",
			},
			wantErr: false,
		},
		{
			name: "множественные пробелы внутри",
			args: args{
				value: "высшее   юридическое    образование",
			},
			want: Education{
				value: "высшее юридическое образование",
			},
			wantErr: false,
		},
		{
			name: "табуляции и переводы строк внутри",
			args: args{
				value: "высшее\tобразование\nуниверситет",
			},
			want: Education{
				value: "высшее образование университет",
			},
			wantErr: false,
		},
		{
			name: "смешанные пробельные символы по краям и внутри",
			args: args{
				value: "  высшее \t образование \n университет  ",
			},
			want: Education{
				value: "высшее образование университет",
			},
			wantErr: false,
		},
		{
			name: "unicode строка",
			args: args{
				value: "  Привет,   мир!  ",
			},
			want: Education{
				value: "Привет, мир!",
			},
			wantErr: false,
		},
		{
			name: "unicode ровно лимит",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit),
			},
			want: Education{
				value: strings.Repeat("Я", ValueCharsLimit),
			},
			wantErr: false,
		},
		{
			name: "unicode за лимитом",
			args: args{
				value: strings.Repeat("Я", ValueCharsLimit+1),
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "эмодзи в пределах лимита",
			args: args{
				value: "  🎓  образование  🎓  ",
			},
			want: Education{
				value: "🎓 образование 🎓",
			},
			wantErr: false,
		},
		{
			name: "эмодзи за лимитом",
			args: args{
				value: strings.Repeat("🎓", ValueCharsLimit+1),
			},
			want:    Education{},
			wantErr: true,
		},
		{
			name: "реалистичные сведения об образовании",
			args: args{
				value: "  Высшее,   Московский государственный университет  им. М.В. Ломоносова, 2005-2010  ",
			},
			want: Education{
				value: "Высшее, Московский государственный университет им. М.В. Ломоносова, 2005-2010",
			},
			wantErr: false,
		},
		{
			name: "carriage return внутри",
			args: args{
				value: "образование\rуниверситет",
			},
			want: Education{
				value: "образование университет",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.value)
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

func TestEducation_Value(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "пустая строка",
			fields: fields{
				value: "",
			},
			want: "",
		},
		{
			name: "один ASCII символ",
			fields: fields{
				value: "a",
			},
			want: "a",
		},
		{
			name: "короткая строка",
			fields: fields{
				value: strings.Repeat("a", 20),
			},
			want: strings.Repeat("a", 20),
		},
		{
			name: "длинная строка до лимита",
			fields: fields{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: strings.Repeat("a", ValueCharsLimit),
		},
		{
			name: "строка с пробелами",
			fields: fields{
				value: "высшее образование",
			},
			want: "высшее образование",
		},
		{
			name: "строка с табуляцией",
			fields: fields{
				value: "высшее\tобразование",
			},
			want: "высшее\tобразование",
		},
		{
			name: "unicode строка",
			fields: fields{
				value: "Привет, мир!",
			},
			want: "Привет, мир!",
		},
		{
			name: "строка с эмодзи",
			fields: fields{
				value: "🎓 образование 🎓",
			},
			want: "🎓 образование 🎓",
		},
		{
			name: "нулевые символы",
			fields: fields{
				value: "\x00\x00\x00",
			},
			want: "\x00\x00\x00",
		},
		{
			name: "символы нулевой ширины",
			fields: fields{
				value: "\u200B\u200B\u200B",
			},
			want: "\u200B\u200B\u200B",
		},
		{
			name: "смешанные пробельные символы",
			fields: fields{
				value: " \t\n\r ",
			},
			want: " \t\n\r ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Education{
				value: tt.fields.value,
			}
			if got := e.Value(); got != tt.want {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEducation_IsZero(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "пустая строка",
			fields: fields{
				value: "",
			},
			want: true,
		},
		{
			name: "один ASCII символ",
			fields: fields{
				value: "a",
			},
			want: false,
		},
		{
			name: "строка из пробелов",
			fields: fields{
				value: "   ",
			},
			want: false,
		},
		{
			name: "строка из табуляций",
			fields: fields{
				value: "\t\t\t",
			},
			want: false,
		},
		{
			name: "строка из переводов строк",
			fields: fields{
				value: "\n\n\n",
			},
			want: false,
		},
		{
			name: "unicode строка",
			fields: fields{
				value: "Привет, мир!",
			},
			want: false,
		},
		{
			name: "строка с эмодзи",
			fields: fields{
				value: "🎓 образование 🎓",
			},
			want: false,
		},
		{
			name: "нулевые символы",
			fields: fields{
				value: "\x00\x00\x00",
			},
			want: false,
		},
		{
			name: "символы нулевой ширины",
			fields: fields{
				value: "\u200B\u200B\u200B",
			},
			want: false,
		},
		{
			name: "длинная строка до лимита",
			fields: fields{
				value: strings.Repeat("a", ValueCharsLimit),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := Education{
				value: tt.fields.value,
			}
			if got := e.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}
