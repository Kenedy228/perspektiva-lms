package passhash

import "testing"

func Test_normalizeHash(t *testing.T) {
	type args struct {
		hash string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "удалить незначащие пробелы в начале и конце строки",
			args: args{
				hash: " hash123  ",
			},
			want: "hash123",
		},
		{
			name: "оставить как есть",
			args: args{
				hash: "hash123",
			},
			want: "hash123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeHash(tt.args.hash); got != tt.want {
				t.Errorf("normalizeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}
