package sequence

import (
	"reflect"
	"testing"

	question2 "gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/option"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
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

func mustText(t *testing.T, v string) text.Text {
	t.Helper()

	got, err := text.New(v)
	if err != nil {
		t.Fatalf("text.New() error = %v", err)
	}

	return got
}

func mustOption(t *testing.T, v string) option.Option {
	t.Helper()

	got, err := option.New(mustText(t, v))
	if err != nil {
		t.Fatalf("option.New() error = %v", err)
	}

	return got
}

func makeOptions(t *testing.T, n int) []option.Option {
	t.Helper()

	options := make([]option.Option, 0, n)
	for i := 0; i < n; i++ {
		options = append(options, mustOption(t, uuid.NewString()))
	}

	return options
}

func TestNew(t *testing.T) {
	type args struct {
		t       title.Title
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
				t:       mustTitle(t, "sequence"),
				options: makeOptions(t, MinOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "error when options less than min",
			args: args{
				t:       mustTitle(t, "sequence"),
				options: makeOptions(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "error when options more than max",
			args: args{
				t:       mustTitle(t, "sequence"),
				options: makeOptions(t, MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "error when contains empty option",
			args: args{
				t: mustTitle(t, "sequence"),
				options: []option.Option{
					mustOption(t, "a"),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "error when contains duplicate ids",
			args: args{
				t: mustTitle(t, "sequence"),
				options: func() []option.Option {
					o := mustOption(t, "a")
					return []option.Option{o, o}
				}(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got == nil {
				t.Fatal("New() got nil, want non-nil")
			}
			if got.Type() != question2.TypeSequence {
				t.Errorf("Type() = %v, want %v", got.Type(), question2.TypeSequence)
			}
			if !reflect.DeepEqual(got.Options(), tt.args.options) {
				t.Errorf("Options() = %v, want %v", got.Options(), tt.args.options)
			}
		})
	}
}

func TestQuestion_ChangeOptions(t *testing.T) {
	initialOptions := makeOptions(t, MinOptionsCount)
	newOptions := makeOptions(t, MinOptionsCount+1)

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
			name: "changes options on valid input",
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
			name: "does not change options on invalid input",
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
	type fields struct {
		Base    *base.Base
		options []option.Option
	}
	tests := []struct {
		name   string
		fields fields
		want   question2.Question
	}{
		{
			name: "returns cloned question",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, MinOptionsCount),
			},
			want: &Question{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, MinOptionsCount),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:    tt.fields.Base,
				options: tt.fields.options,
			}
			got, ok := q.Clone().(*Question)
			if !ok {
				t.Fatalf("Clone() type = %T, want *Question", q.Clone())
			}
			if got == q {
				t.Errorf("Clone() returned same pointer")
			}
			if got.Base == q.Base {
				t.Errorf("Clone() Base pointer is same, want different")
			}
			if !reflect.DeepEqual(got.Options(), q.Options()) {
				t.Errorf("Clone().Options() = %v, want %v", got.Options(), q.Options())
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
				options: makeOptions(t, MinOptionsCount),
			},
			want: question2.TypeSequence.DefaultInstruction(),
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
	options := makeOptions(t, MinOptionsCount)

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
			name: "returns copy of options",
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
			got := q.Options()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Options() = %v, want %v", got, tt.want)
			}
			if len(got) > 0 && len(tt.fields.options) > 0 && &got[0] == &tt.fields.options[0] {
				t.Errorf("Options() must return cloned slice")
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
		want   question2.Type
	}{
		{
			name: "returns sequence type",
			fields: fields{
				Base:    mustBase(t, "q"),
				options: makeOptions(t, MinOptionsCount),
			},
			want: question2.TypeSequence,
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
				options: makeOptions(t, MinOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "less than min",
			args: args{
				options: makeOptions(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "more than max",
			args: args{
				options: makeOptions(t, MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "contains empty option",
			args: args{
				options: []option.Option{
					mustOption(t, "a"),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "contains duplicate ids",
			args: args{
				options: func() []option.Option {
					o := mustOption(t, "a")
					return []option.Option{o, o}
				}(),
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
				options: makeOptions(t, MinOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "has empty option",
			args: args{
				options: []option.Option{
					mustOption(t, "a"),
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty slice allowed",
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
			name: "exact min",
			args: args{
				options: makeOptions(t, MinOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "exact max",
			args: args{
				options: makeOptions(t, MaxOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "less than min",
			args: args{
				options: makeOptions(t, MinOptionsCount-1),
			},
			wantErr: true,
		},
		{
			name: "more than max",
			args: args{
				options: makeOptions(t, MaxOptionsCount+1),
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
			name: "no duplicates",
			args: args{
				options: makeOptions(t, MinOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "has duplicates",
			args: args{
				options: func() []option.Option {
					o := mustOption(t, "a")
					return []option.Option{o, o}
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
