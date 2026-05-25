package education

import "testing"

func Test_normalizeValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "пустая строка",
			args: args{
				value: "",
			},
			want: "",
		},
		{
			name: "только пробелы",
			args: args{
				value: "   ",
			},
			want: "",
		},
		{
			name: "только табуляции",
			args: args{
				value: "\t\t\t",
			},
			want: "",
		},
		{
			name: "только переводы строк",
			args: args{
				value: "\n\n\n",
			},
			want: "",
		},
		{
			name: "смешанные пробельные символы без текста",
			args: args{
				value: " \t \n  \r",
			},
			want: "",
		},
		{
			name: "ведущие и замыкающие пробелы",
			args: args{
				value: "  высшее  ",
			},
			want: "высшее",
		},
		{
			name: "ведущие и замыкающие табуляции",
			args: args{
				value: "\t\tвысшее\t\t",
			},
			want: "высшее",
		},
		{
			name: "ведущие и замыкающие переводы строк",
			args: args{
				value: "\n\nвысшее\n\n",
			},
			want: "высшее",
		},
		{
			name: "дубликаты пробелов внутри",
			args: args{
				value: "высшее  юридическое   образование",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "табуляции внутри",
			args: args{
				value: "высшее\tюридическое\t\tобразование",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "переводы строк внутри",
			args: args{
				value: "высшее\nюридическое\n\nобразование",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "смешанные пробельные символы внутри и по краям",
			args: args{
				value: "  высшее \t юридическое \n  образование  ",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "одинарный пробел не меняется",
			args: args{
				value: "высшее юридическое",
			},
			want: "высшее юридическое",
		},
		{
			name: "строка без пробелов",
			args: args{
				value: "образование",
			},
			want: "образование",
		},
		{
			name: "один символ с пробелами",
			args: args{
				value: "  a  ",
			},
			want: "a",
		},
		{
			name: "unicode с пробелами",
			args: args{
				value: "  Привет,   мир!  ",
			},
			want: "Привет, мир!",
		},
		{
			name: "эмодзи с пробелами",
			args: args{
				value: "  🎓  образование  🎓  ",
			},
			want: "🎓 образование 🎓",
		},
		{
			name: "carriage return и form feed",
			args: args{
				value: "высшее\rюридическое\fобразование",
			},
			want: "высшее юридическое образование",
		},
		{
			name: "вертикальная табуляция",
			args: args{
				value: "высшее\vюридическое",
			},
			want: "высшее юридическое",
		},
		{
			name: "только пробел в начале",
			args: args{
				value: " высшее",
			},
			want: "высшее",
		},
		{
			name: "только пробел в конце",
			args: args{
				value: "высшее ",
			},
			want: "высшее",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeValue(tt.args.value); got != tt.want {
				t.Errorf("normalizeValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
