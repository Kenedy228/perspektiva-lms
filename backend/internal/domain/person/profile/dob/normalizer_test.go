package dob

import (
	"reflect"
	"testing"
	"time"
)

func Test_normalize(t *testing.T) {
	type args struct {
		v time.Time
	}
	tests := []struct {
		name string
		args args
		want time.Time
	}{
		{
			name: "уже нормализованное значение",
			args: args{
				v: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			want: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "не нормализованное значение",
			args: args{
				v: time.Date(2000, 5, 10, 15, 4, 5, 123, time.FixedZone("MSK", 3*3600)),
			},
			want: time.Date(2000, 5, 10, 0, 0, 0, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalize(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("normalize() = %v, want %v", got, tt.want)
			}
		})
	}
}
