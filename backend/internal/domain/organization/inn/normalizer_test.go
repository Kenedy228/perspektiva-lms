package inn

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
			name: "удаляет все пробелы",
			args: args{
				value: " 11 22    33 44      ",
			},
			want: "11223344",
		},
		{
			name: "не изменяет значение без пробелов",
			args: args{
				value: "11223344",
			},
			want: "11223344",
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
