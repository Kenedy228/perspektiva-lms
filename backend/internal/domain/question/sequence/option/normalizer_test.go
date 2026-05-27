package option

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
				value: "   value    ",
			},
			want: "value",
		},
		{
			name: "заменяет последовательности пробелов внутри предложения на одиночный пробел",
			args: args{
				value: "va      lue",
			},
			want: "va lue",
		},
		{
			name: "удаляет незначащие пробелы и заменяет последовательность пробелов на одиночный пробел",
			args: args{
				value: "   va     lue    ",
			},
			want: "va lue",
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
