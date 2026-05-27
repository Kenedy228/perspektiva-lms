package sequence

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	b := baseFixture(t)

	type args struct {
		b       *base.Base
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Question
		wantErr bool
	}{
		{
			name: "опций меньше порогового значения",
			args: args{
				b:       baseFixture(t),
				options: optionsSliceFixture(t, MinOptionsCount-1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "опций больше порогового значения",
			args: args{
				b:       baseFixture(t),
				options: optionsSliceFixture(t, MaxOptionsCount+1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "содержит пустую опцию",
			args: args{
				b: baseFixture(t),
				options: func() []option.Option {
					res := make([]option.Option, 0, MinOptionsCount+1)
					for range MinOptionsCount {
						v := option.Option{}
						res = append(res, v)
					}
					return res
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет базы",
			args: args{
				b:       nil,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет опций и базы",
			args: args{
				b:       nil,
				options: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "валидное число опций и база",
			args: args{
				b:       b,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			want: &Question{
				Base:    b,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.b, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				assert.NotEqual(t, uuid.Nil, got.ID())
				assert.Equal(t, tt.want.Base, got.Base)
				assert.Equal(t, tt.want.Options(), got.Options())
			}
		})
	}
}

func TestQuestion_ChangeOptions(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "опций меньше порогового значения",
			fields: fields{
				Base:    baseFixture(t),
				options: []option.Option{},
			},
			args: args{
				options: optionsSliceFixture(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "опций больше порогового значения",
			fields: fields{
				Base:    baseFixture(t),
				options: []option.Option{},
			},
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "содержит пустую опцию",
			fields: fields{
				Base:    baseFixture(t),
				options: []option.Option{},
			},
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
			fields: fields{
				Base:    baseFixture(t),
				options: []option.Option{},
			},
			args: args{
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			if err := q.ChangeOptions(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("ChangeOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestion_Clone(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "возвращает копию вопроса",
			fields: fields{
				Base:    &base.Base{},
				options: []option.Option{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}

			got := q.Clone()

			assert.Equal(t, q.ID(), got.ID())
			assert.Equal(t, q.Title(), got.Title())
			assert.Equal(t, q.Instruction(), got.Instruction())
			assert.Equal(t, q.Type(), got.Type())
		})
	}
}

func TestQuestion_Instruction(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращает инструкцию TypeSequence",
			fields: fields{
				Base:    &base.Base{},
				options: []option.Option{},
			},
			want: question.TypeSequence.DefaultInstruction(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			if got := q.Instruction(); got != tt.want {
				t.Errorf("Instruction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Options(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
		want   []option.Option
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				Base:    &base.Base{},
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			want: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			if got := q.Options(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Type(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
		want   question.Type
	}{
		{
			name: "возвращает TypeSequence",
			fields: fields{
				Base:    &base.Base{},
				options: []option.Option{},
			},
			want: question.TypeSequence,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			if got := q.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	b := baseFixture(t)

	type args struct {
		b       *base.Base
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		want    *Question
		wantErr bool
	}{
		{
			name: "опций меньше порогового значения",
			args: args{
				b:       baseFixture(t),
				options: optionsSliceFixture(t, MinOptionsCount-1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "опций больше порогового значения",
			args: args{
				b:       baseFixture(t),
				options: optionsSliceFixture(t, MaxOptionsCount+1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "содержит пустую опцию",
			args: args{
				b: baseFixture(t),
				options: func() []option.Option {
					res := make([]option.Option, 0, MinOptionsCount+1)
					for range MinOptionsCount {
						v := option.Option{}
						res = append(res, v)
					}
					return res
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет базы",
			args: args{
				b:       nil,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет опций и базы",
			args: args{
				b:       nil,
				options: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "валидное число опций и база",
			args: args{
				b:       b,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			want: &Question{
				Base:    b,
				options: optionsSliceFixture(t, MaxOptionsCount-MinOptionsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.b, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Restore() got = %v, want %v", got, tt.want)
			}
		})
	}
}
