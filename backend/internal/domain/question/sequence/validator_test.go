package sequence

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
)

func Test_validateBase(t *testing.T) {
	type args struct {
		b *base.Base
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустая база",
			args: args{
				b: nil,
			},
			wantErr: true,
		},
		{
			name: "непустая база",
			args: args{
				b: baseFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBase(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("validateBase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateBaseRequired(t *testing.T) {
	type args struct {
		b *base.Base
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "пустая база",
			args: args{
				b: nil,
			},
			wantErr: true,
		},
		{
			name: "непустая база",
			args: args{
				b: baseFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBaseRequired(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("validateBaseRequired() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptions(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "опций меньше порогового значения",
			args: args{
				options: optionsSliceFixture(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "опций больше порогового значения",
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "содержит пустую опцию",
			args: args{
				options: func() []option.Option {
					res := make([]option.Option, 0, MinOptionsCount+1)
					for range MinOptionsCount {
						v := option.Option{}
						res = append(res, v)
					}
					return res
				}(),
			},
			wantErr: true,
		},
		{
			name: "валидное число опций",
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptions(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("validateOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionsContainsEmpty(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "содержит пустую опцию",
			args: args{
				options: func() []option.Option {
					res := make([]option.Option, 0, MinOptionsCount+1)
					for range MinOptionsCount {
						v := option.Option{}
						res = append(res, v)
					}
					return res
				}(),
			},
			wantErr: true,
		},
		{
			name: "не содержит пустых опций",
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionsContainsEmpty(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionsContainsEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionsCount(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "опций меньше порогового значения",
			args: args{
				options: optionsSliceFixture(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "опций больше порогового значения",
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "валидное число опций",
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionsCount(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
