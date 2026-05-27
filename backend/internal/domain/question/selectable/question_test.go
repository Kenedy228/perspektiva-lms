package selectable

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
	"github.com/google/uuid"
)

func mustTitle(t *testing.T, v string) title.Title {
	t.Helper()

	got, err := title.New(v)
	if err != nil {
		t.Fatalf("title.New() error = %v", err)
	}

	return got
}

func mustBase(t *testing.T, v string) *base.Base {
	t.Helper()

	got, err := base.New(mustTitle(t, v))
	if err != nil {
		t.Fatalf("base.New() error = %v", err)
	}

	return got
}

func mustOption(t *testing.T, v string, isCorrect bool) option.Option {
	t.Helper()

	got, err := option.New(v, isCorrect)
	if err != nil {
		t.Fatalf("option.New() error = %v", err)
	}

	return got
}

func makeOptions(t *testing.T, n int, correctCount int) []option.Option {
	t.Helper()

	options := make([]option.Option, 0, n)
	for i := 0; i < n; i++ {
		options = append(options, mustOption(t, uuid.NewString(), i < correctCount))
	}

	return options
}

func TestNew(t *testing.T) {
	validOptions := makeOptions(t, MinOptionsCount, 1)

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
			name: "creates question with valid options",
			args: args{
				b:       mustBase(t, "selectable"),
				options: validOptions,
			},
			wantErr: false,
		},
		{
			name: "returns error when options less than min",
			args: args{
				b:       mustBase(t, "selectable"),
				options: makeOptions(t, MinOptionsCount-1, 1),
			},
			wantErr: true,
		},
		{
			name: "returns error when no correct options",
			args: args{
				b:       mustBase(t, "selectable"),
				options: makeOptions(t, MinOptionsCount, 0),
			},
			wantErr: true,
		},
		{
			name: "returns error when contains empty option",
			args: args{
				b: mustBase(t, "selectable"),
				options: []option.Option{
					mustOption(t, "a", true),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "returns error when contains duplicate option id",
			args: args{
				b: mustBase(t, "selectable"),
				options: func() []option.Option {
					opt := mustOption(t, "a", true)
					return []option.Option{opt, opt}
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.b, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Fatal("New() got nil, want non-nil question")
			}
			if got.Type() != question.TypeSelectable {
				t.Errorf("Type() = %v, want %v", got.Type(), question.TypeSelectable)
			}
			if !reflect.DeepEqual(got.Options(), tt.args.options) {
				t.Errorf("Options() = %v, want %v", got.Options(), tt.args.options)
			}
		})
	}
}

func TestQuestion_ChangeOptions(t *testing.T) {
	initialOptions := makeOptions(t, MinOptionsCount, 1)
	newOptions := makeOptions(t, MinOptionsCount+1, 2)

	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantErr     bool
		wantOptions []option.Option
	}{
		{
			name: "changes options with valid input",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: initialOptions,
			},
			args: args{
				options: newOptions,
			},
			wantErr:     false,
			wantOptions: newOptions,
		},
		{
			name: "does not change options with invalid input",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: initialOptions,
			},
			args: args{
				options: []option.Option{{}},
			},
			wantErr:     true,
			wantOptions: initialOptions,
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
			if !reflect.DeepEqual(q.Options(), tt.wantOptions) {
				t.Errorf("ChangeOptions() options = %v, want %v", q.Options(), tt.wantOptions)
			}
		})
	}
}

func TestQuestion_Clone(t *testing.T) {
	options := makeOptions(t, MinOptionsCount, 1)
	q := &Question{
		Base:    mustBase(t, "q"),
		options: options,
	}

	got, ok := q.Clone().(*Question)
	if !ok {
		t.Fatalf("Clone() type = %T, want *Question", q.Clone())
	}
	if got == q {
		t.Errorf("Clone() returned same pointer")
	}
	if got.Base == q.Base {
		t.Errorf("Clone() Base pointer = same, want different")
	}
	if !reflect.DeepEqual(got.Options(), q.Options()) {
		t.Errorf("Clone() Options() = %v, want %v", got.Options(), q.Options())
	}
}

func TestQuestion_CorrectOptionsCount(t *testing.T) {
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "returns count of correct options",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, 4, 2),
			},
			want: 2,
		},
		{
			name: "returns zero for empty slice",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: []option.Option{},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			if got := q.CorrectOptionsCount(); got != tt.want {
				t.Errorf("CorrectOptionsCount() = %v, want %v", got, tt.want)
			}
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
			name: "returns default instruction",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, MinOptionsCount, 1),
			},
			want: question.TypeSelectable.DefaultInstruction(),
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
	options := makeOptions(t, MinOptionsCount, 1)

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
			name: "returns options",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: options,
			},
			want: options,
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
			name: "returns selectable type",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, MinOptionsCount, 1),
			},
			want: question.TypeSelectable,
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

func Test_countCorrectOptions(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "counts correct options",
			args: args{
				options: makeOptions(t, 5, 3),
			},
			want: 3,
		},
		{
			name: "empty slice",
			args: args{
				options: []option.Option{},
			},
			want: 0,
		},
		{
			name: "no correct options",
			args: args{
				options: makeOptions(t, 5, 0),
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countCorrectOptions(tt.args.options); got != tt.want {
				t.Errorf("countCorrectOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateCorrectOptionsCount(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid correct options count",
			args: args{
				options: makeOptions(t, MinOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "no correct options",
			args: args{
				options: makeOptions(t, MinOptionsCount, 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateCorrectOptionsCount(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("validateCorrectOptionsCount() error = %v, wantErr %v", err, tt.wantErr)
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
			name: "valid options",
			args: args{
				options: makeOptions(t, MinOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "less than min count",
			args: args{
				options: makeOptions(t, MinOptionsCount-1, 1),
			},
			wantErr: true,
		},
		{
			name: "contains empty option",
			args: args{
				options: []option.Option{
					mustOption(t, "a", true),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "contains duplicate options",
			args: args{
				options: func() []option.Option {
					opt := mustOption(t, "a", true)
					return []option.Option{opt, opt}
				}(),
			},
			wantErr: true,
		},
		{
			name: "no correct options",
			args: args{
				options: makeOptions(t, MinOptionsCount, 0),
			},
			wantErr: true,
		},
		{
			name: "more than max count",
			args: args{
				options: makeOptions(t, MaxOptionsCount+1, 1),
			},
			wantErr: true,
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
			name: "no empty options",
			args: args{
				options: makeOptions(t, MinOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "contains empty option",
			args: args{
				options: []option.Option{
					mustOption(t, "a", true),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				options: []option.Option{},
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
			name: "exact min count",
			args: args{
				options: makeOptions(t, MinOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "exact max count",
			args: args{
				options: makeOptions(t, MaxOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "less than min count",
			args: args{
				options: makeOptions(t, MinOptionsCount-1, 1),
			},
			wantErr: true,
		},
		{
			name: "more than max count",
			args: args{
				options: makeOptions(t, MaxOptionsCount+1, 1),
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				options: []option.Option{},
			},
			wantErr: true,
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

func Test_validateOptionsDuplicates(t *testing.T) {
	type args struct {
		options []option.Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unique options",
			args: args{
				options: makeOptions(t, MinOptionsCount, 1),
			},
			wantErr: false,
		},
		{
			name: "duplicate option ids",
			args: args{
				options: func() []option.Option {
					opt := mustOption(t, "same", true)
					return []option.Option{opt, opt}
				}(),
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				options: []option.Option{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionsDuplicates(tt.args.options); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionsDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
