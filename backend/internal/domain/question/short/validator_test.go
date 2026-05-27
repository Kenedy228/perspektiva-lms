package short

import (
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/short/variant"
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
			name: "база nil",
			args: args{
				b: nil,
			},
			wantErr: true,
		},
		{
			name: "база не nil",
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

func Test_validateRequiredBase(t *testing.T) {
	type args struct {
		b *base.Base
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "база nil",
			args: args{
				b: nil,
			},
			wantErr: true,
		},
		{
			name: "база не nil",
			args: args{
				b: baseFixture(t),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateRequiredBase(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("validateRequiredBase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateVariants(t *testing.T) {
	type args struct {
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нет вариантов",
			args: args{
				variants: []variant.Variant{},
			},
			wantErr: true,
		},
		{
			name: "вариантов меньше необходимого количества",
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount-1),
			},
			wantErr: true,
		},
		{
			name: "вариантов больше необходимого количества",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount+1),
			},
			wantErr: true,
		},
		{
			name: "вариантов ровно минимум",
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "вариантов ровно максимум",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "варианты между минимумом и максимумом",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount-MinVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "содержит пустой вариант",
			args: args{
				variants: func() []variant.Variant {
					res := make([]variant.Variant, 0, MinVariantsCount+1)
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
			if err := validateVariants(tt.args.variants); (err != nil) != tt.wantErr {
				t.Errorf("validateVariants() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateVariantsContainsNonInitialized(t *testing.T) {
	type args struct {
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "содержит пустой вариант",
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
		{
			name: "не содержит пустых вариантов",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount-MinVariantsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVariantsContainsNonInitialized(tt.args.variants); (err != nil) != tt.wantErr {
				t.Errorf("validateVariantsContainsNonInitialized() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateVariantsCount(t *testing.T) {
	type args struct {
		variants []variant.Variant
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "нет вариантов",
			args: args{
				variants: []variant.Variant{},
			},
			wantErr: true,
		},
		{
			name: "вариантов меньше необходимого количества",
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount-1),
			},
			wantErr: true,
		},
		{
			name: "вариантов больше необходимого количества",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount+1),
			},
			wantErr: true,
		},
		{
			name: "вариантов ровно минимум",
			args: args{
				variants: variantsSliceFixture(t, MinVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "вариантов ровно максимум",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount),
			},
			wantErr: false,
		},
		{
			name: "варианты между минимумом и максимумом",
			args: args{
				variants: variantsSliceFixture(t, MaxVariantsCount-MinVariantsCount),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateVariantsCount(tt.args.variants); (err != nil) != tt.wantErr {
				t.Errorf("validateVariantsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
