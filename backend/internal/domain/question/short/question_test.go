package short

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type args struct {
		b        *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		args    args
		want    *Question
		wantErr bool
	}{
		{
			name: "нет вариантов",
			args: args{
				b:        baseFixture(t),
				variants: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет базы",
			args: args{
				b:        nil,
				variants: []variant.Variant{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "вариантов меньше необходимого количества",
			args: args{
				b:        baseFixture(t),
				variants: variantsSliceFixture(t, MinVariantsCount-1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "вариантов больше необходимого количества",
			args: args{
				b:        baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount+1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "содержит пустой вариант",
			args: args{
				b: baseFixture(t),
				variants: func() []variant.Variant {
					res := make([]variant.Variant, 0, MaxVariantsCount+1)
					for range MinVariantsCount {
						v := variant.Variant{}
						res = append(res, v)
					}
					return res
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "есть база и варианты",
			args: args{
				b:        baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			want: &Question{
				Base:     baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.b, tt.args.variants)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want != nil {
				assert.NotNil(t, got.Base)
				assert.NotNil(t, got.variants)
				assert.Equal(t, len(tt.want.variants), len(got.variants))
			}
		})
	}
}

func TestQuestion_ChangeVariants(t *testing.T) {
	type fields struct {
		Base     *base.Base
		variants []variant.Variant
	}
	type args struct {
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "нет вариантов",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: nil,
			},
			wantErr: true,
		},
		{
			name: "пустые варианты",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: []variant.Variant{},
			},
			wantErr: true,
		},
		{
			name: "вариантов меньше необходимого количества",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount-1),
			},
			wantErr: true,
		},
		{
			name: "вариантов больше необходимого количества",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount+1),
			},
			wantErr: true,
		},
		{
			name: "вариантов ровно минимум",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "вариантов ровно максимум",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "варианты между минимумом и максимумом",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount-MinVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "содержит пустой вариант",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			args: args{
				variants: func() []variant.Variant {
					res := make([]variant.Variant, 0, MaxVariantsCount+1)
					for range MinVariantsCount {
						v := variant.Variant{}
						res = append(res, v)
					}
					return res
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:     tt.fields.Base,
				variants: tt.fields.variants,
			}
			if err := q.ChangeVariants(tt.args.variants); (err != nil) != tt.wantErr {
				t.Errorf("ChangeVariants() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQuestion_Clone(t *testing.T) {
	type fields struct {
		Base     *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "клонирует значение как есть",
			fields: fields{
				Base:     baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:     tt.fields.Base,
				variants: tt.fields.variants,
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
		Base     *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "возвращается инструкция типа Short",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			want: question.TypeShort.DefaultInstruction(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:     tt.fields.Base,
				variants: tt.fields.variants,
			}
			if got := q.Instruction(); got != tt.want {
				t.Errorf("Instruction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Type(t *testing.T) {
	type fields struct {
		Base     *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name   string
		fields fields
		want   question.Type
	}{
		{
			name: "возвращается тип Short",
			fields: fields{
				Base:     nil,
				variants: nil,
			},
			want: question.TypeShort,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:     tt.fields.Base,
				variants: tt.fields.variants,
			}
			if got := q.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Variants(t *testing.T) {
	type fields struct {
		Base     *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name   string
		fields fields
		want   []variant.Variant
	}{
		{
			name: "возвращает значение как есть",
			fields: fields{
				Base:     nil,
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			want: variantsSliceFixture(t, MaxVariantsCount-1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:     tt.fields.Base,
				variants: tt.fields.variants,
			}
			if got := q.Variants(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variants() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRestore(t *testing.T) {
	type args struct {
		b        *base.Base
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		args    args
		want    *Question
		wantErr bool
	}{
		{
			name: "валидные база и варианты",
			args: args{
				b:        baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			want: &Question{
				Base:     baseFixture(t),
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			wantErr: false,
		},
		{
			name: "нет базы",
			args: args{
				b:        nil,
				variants: variantsSliceFixture(t, MaxVariantsCount-1),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "нет вариантов",
			args: args{
				b:        baseFixture(t),
				variants: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.args.b, tt.args.variants)
			if (err != nil) != tt.wantErr {
				t.Errorf("Restore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				assert.NotNil(t, got.Base)
				assert.NotNil(t, got.variants)
				assert.Equal(t, len(tt.want.variants), len(got.variants))
			}
		})
	}
}
