package jobtitle

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
			name: "удаляет незначащие пробелы",
			args: args{
				value: " монтер ",
			},
			want: "монтер",
		},
		{
			name: "удаляет дубликаты пробелов внутри",
			args: args{
				value: " монтер  механик  \tэлектро",
			},
			want: "монтер механик электро",
		},
		{
			name: "не меняет валидную строку",
			args: args{
				value: "механик",
			},
			want: "механик",
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
