package name

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
			name: "все последовательность из нескольких пробелов на одиночный пробел",
			args: args{
				value: "    ООО    Ромашка   ",
			},
			want: "ООО Ромашка",
		},
		{
			name: "не меняет нормализованное значение",
			args: args{
				value: "ООО Ромашка",
			},
			want: "ООО Ромашка",
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
