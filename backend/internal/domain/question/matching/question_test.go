package matching

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	pair2 "gitflic.ru/lms/backend/internal/domain/question/matching/pair"
	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
)

func mustText(t *testing.T, value string) text.Text {
	t.Helper()

	got, err := text.New(value)
	if err != nil {
		t.Fatalf("text.New() error = %v", err)
	}

	return got
}

func mustTitle(t *testing.T, value string) title.Title {
	t.Helper()

	got, err := title.New(value)
	if err != nil {
		t.Fatalf("title.New() error = %v", err)
	}

	return got
}

func mustBase(t *testing.T, value string) *base.Base {
	t.Helper()

	got, err := base.New(mustTitle(t, value))
	if err != nil {
		t.Fatalf("base.New() error = %v", err)
	}

	return got
}

func mustPair(t *testing.T, promptValue, matchValue string) pair2.Pair {
	t.Helper()

	prompt, err := pair2.NewPrompt(mustText(t, promptValue))
	if err != nil {
		t.Fatalf("pair.NewPrompt() error = %v", err)
	}

	match, err := pair2.NewMatch(mustText(t, matchValue))
	if err != nil {
		t.Fatalf("pair.NewMatch() error = %v", err)
	}

	got, err := pair2.New(prompt, match)
	if err != nil {
		t.Fatalf("pair.New() error = %v", err)
	}

	return got
}

func makePairs(t *testing.T, n int) []pair2.Pair {
	t.Helper()

	pairs := make([]pair2.Pair, 0, n)
	for i := 0; i < n; i++ {
		pairs = append(pairs, mustPair(t, "prompt", "match"))
	}

	return pairs
}

func TestNew(t *testing.T) {
	validTitle := mustTitle(t, "matching question")
	validPairs := makePairs(t, MinPairs)

	type args struct {
		t     title.Title
		pairs []pair2.Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid question",
			args: args{
				t:     validTitle,
				pairs: validPairs,
			},
			wantErr: false,
		},
		{
			name: "invalid not enough pairs",
			args: args{
				t:     validTitle,
				pairs: makePairs(t, MinPairs-1),
			},
			wantErr: true,
		},
		{
			name: "invalid empty pair inside",
			args: args{
				t: validTitle,
				pairs: []pair2.Pair{
					mustPair(t, "p1", "m1"),
					pair2.Pair{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid too many pairs",
			args: args{
				t:     validTitle,
				pairs: makePairs(t, MaxPairs+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.t, tt.args.pairs)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got == nil {
				t.Fatalf("New() got = nil, want non-nil question")
			}
			if got.Type() != question.TypeMatching {
				t.Errorf("New() Type() = %v, want %v", got.Type(), question.TypeMatching)
			}
			if !reflect.DeepEqual(got.Pairs(), tt.args.pairs) {
				t.Errorf("New() Pairs() = %v, want %v", got.Pairs(), tt.args.pairs)
			}
		})
	}
}

func TestQuestion_ChangePairs(t *testing.T) {
	oldPairs := makePairs(t, MinPairs)
	newPairs := makePairs(t, MinPairs+1)

	type fields struct {
		Base  *base.Base
		pairs []pair2.Pair
	}
	type args struct {
		pairs []pair2.Pair
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   bool
		wantPairs []pair2.Pair
	}{
		{
			name: "change to valid pairs",
			fields: fields{
				Base:  mustBase(t, "question"),
				pairs: oldPairs,
			},
			args: args{
				pairs: newPairs,
			},
			wantErr:   false,
			wantPairs: newPairs,
		},
		{
			name: "reject invalid pairs and keep old state",
			fields: fields{
				Base:  mustBase(t, "question"),
				pairs: oldPairs,
			},
			args: args{
				pairs: []pair2.Pair{pair2.Pair{}},
			},
			wantErr:   true,
			wantPairs: oldPairs,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:  tt.fields.Base,
				pairs: tt.fields.pairs,
			}
			if err := q.ChangePairs(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("ChangePairs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got := q.Pairs(); !reflect.DeepEqual(got, tt.wantPairs) {
				t.Errorf("ChangePairs() pairs = %v, want %v", got, tt.wantPairs)
			}
		})
	}
}

func TestQuestion_Clone(t *testing.T) {
	base1 := mustBase(t, "question")
	pairs1 := makePairs(t, MinPairs)

	type fields struct {
		Base  *base.Base
		pairs []pair2.Pair
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "clone question",
			fields: fields{
				Base:  base1,
				pairs: pairs1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:  tt.fields.Base,
				pairs: tt.fields.pairs,
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
			if !reflect.DeepEqual(got.Pairs(), q.Pairs()) {
				t.Errorf("Clone() Pairs() = %v, want %v", got.Pairs(), q.Pairs())
			}
			if got.Type() != q.Type() {
				t.Errorf("Clone() Type() = %v, want %v", got.Type(), q.Type())
			}
		})
	}
}

func TestQuestion_Instruction(t *testing.T) {
	type fields struct {
		Base  *base.Base
		pairs []pair2.Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "returns default matching instruction",
			fields: fields{
				Base:  mustBase(t, "question"),
				pairs: makePairs(t, MinPairs),
			},
			want: question.TypeMatching.DefaultInstruction(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:  tt.fields.Base,
				pairs: tt.fields.pairs,
			}
			if got := q.Instruction(); got != tt.want {
				t.Errorf("Instruction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Pairs(t *testing.T) {
	pairs1 := makePairs(t, MinPairs)

	type fields struct {
		Base  *base.Base
		pairs []pair2.Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   []pair2.Pair
	}{
		{
			name: "returns pairs",
			fields: fields{
				Base:  mustBase(t, "question"),
				pairs: pairs1,
			},
			want: pairs1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:  tt.fields.Base,
				pairs: tt.fields.pairs,
			}
			if got := q.Pairs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuestion_Type(t *testing.T) {
	type fields struct {
		Base  *base.Base
		pairs []pair2.Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   question.Type
	}{
		{
			name: "returns matching type",
			fields: fields{
				Base:  mustBase(t, "question"),
				pairs: makePairs(t, MinPairs),
			},
			want: question.TypeMatching,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Question{
				Base:  tt.fields.Base,
				pairs: tt.fields.pairs,
			}
			if got := q.Type(); got != tt.want {
				t.Errorf("Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validatePairs(t *testing.T) {
	type args struct {
		pairs []pair2.Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid pairs",
			args: args{
				pairs: makePairs(t, MinPairs),
			},
			wantErr: false,
		},
		{
			name: "less than min",
			args: args{
				pairs: makePairs(t, MinPairs-1),
			},
			wantErr: true,
		},
		{
			name: "contains empty pair",
			args: args{
				pairs: []pair2.Pair{
					mustPair(t, "p1", "m1"),
					pair2.Pair{},
				},
			},
			wantErr: true,
		},
		{
			name: "more than max",
			args: args{
				pairs: makePairs(t, MaxPairs+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairs(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validatePairs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairsContainsEmpty(t *testing.T) {
	type args struct {
		pairs []pair2.Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all pairs filled",
			args: args{
				pairs: makePairs(t, MinPairs),
			},
			wantErr: false,
		},
		{
			name: "contains zero pair",
			args: args{
				pairs: []pair2.Pair{
					mustPair(t, "p1", "m1"),
					pair2.Pair{},
				},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				pairs: []pair2.Pair{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairsContainsEmpty(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validatePairsContainsEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairsCount(t *testing.T) {
	type args struct {
		pairs []pair2.Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exact min",
			args: args{
				pairs: makePairs(t, MinPairs),
			},
			wantErr: false,
		},
		{
			name: "exact max",
			args: args{
				pairs: makePairs(t, MaxPairs),
			},
			wantErr: false,
		},
		{
			name: "below min",
			args: args{
				pairs: makePairs(t, MinPairs-1),
			},
			wantErr: true,
		},
		{
			name: "above max",
			args: args{
				pairs: makePairs(t, MaxPairs+1),
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				pairs: []pair2.Pair{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairsCount(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validatePairsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
