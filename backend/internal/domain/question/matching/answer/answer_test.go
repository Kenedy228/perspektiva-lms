package answer

import (
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"github.com/google/uuid"
)

func mustPromptID(t *testing.T, id uuid.UUID) PromptID {
	t.Helper()

	got, err := NewPromptID(id)
	if err != nil {
		t.Fatalf("NewPromptID() error = %v", err)
	}

	return got
}

func mustMatchID(t *testing.T, id uuid.UUID) MatchID {
	t.Helper()

	got, err := NewMatchID(id)
	if err != nil {
		t.Fatalf("NewMatchID() error = %v", err)
	}

	return got
}

func makeValidPair(t *testing.T) Pair {
	t.Helper()

	return Pair{
		PromptID: mustPromptID(t, uuid.New()),
		MatchID:  mustMatchID(t, uuid.New()),
	}
}

func makeValidPairs(t *testing.T, n int) []Pair {
	t.Helper()

	pairs := make([]Pair, 0, n)
	for i := 0; i < n; i++ {
		pairs = append(pairs, makeValidPair(t))
	}

	return pairs
}

func TestAnswer_AsMap(t *testing.T) {
	promptID1 := mustPromptID(t, uuid.New())
	matchID1 := mustMatchID(t, uuid.New())
	promptID2 := mustPromptID(t, uuid.New())
	matchID2 := mustMatchID(t, uuid.New())

	type fields struct {
		pairs []Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   map[PromptID]MatchID
	}{
		{
			name: "returns empty map for empty pairs",
			fields: fields{
				pairs: []Pair{},
			},
			want: map[PromptID]MatchID{},
		},
		{
			name: "returns map with all pairs",
			fields: fields{
				pairs: []Pair{
					{PromptID: promptID1, MatchID: matchID1},
					{PromptID: promptID2, MatchID: matchID2},
				},
			},
			want: map[PromptID]MatchID{
				promptID1: matchID1,
				promptID2: matchID2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				pairs: tt.fields.pairs,
			}
			if got := a.AsMap(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_Clone(t *testing.T) {
	promptID1 := mustPromptID(t, uuid.New())
	matchID1 := mustMatchID(t, uuid.New())
	promptID2 := mustPromptID(t, uuid.New())
	matchID2 := mustMatchID(t, uuid.New())

	type fields struct {
		pairs []Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   question.Answer
	}{
		{
			name: "clone empty answer",
			fields: fields{
				pairs: []Pair{},
			},
			want: Answer{
				pairs: []Pair{},
			},
		},
		{
			name: "clone answer with pairs",
			fields: fields{
				pairs: []Pair{
					{PromptID: promptID1, MatchID: matchID1},
					{PromptID: promptID2, MatchID: matchID2},
				},
			},
			want: Answer{
				pairs: []Pair{
					{PromptID: promptID1, MatchID: matchID1},
					{PromptID: promptID2, MatchID: matchID2},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				pairs: tt.fields.pairs,
			}
			if got := a.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_IsEmpty(t *testing.T) {
	type fields struct {
		pairs []Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty answer",
			fields: fields{
				pairs: []Pair{},
			},
			want: true,
		},
		{
			name: "non empty answer",
			fields: fields{
				pairs: []Pair{
					makeValidPair(t),
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				pairs: tt.fields.pairs,
			}
			if got := a.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_Pairs(t *testing.T) {
	pair1 := makeValidPair(t)
	pair2 := makeValidPair(t)

	type fields struct {
		pairs []Pair
	}
	tests := []struct {
		name   string
		fields fields
		want   []Pair
	}{
		{
			name: "returns empty slice",
			fields: fields{
				pairs: []Pair{},
			},
			want: []Pair{},
		},
		{
			name: "returns cloned pairs",
			fields: fields{
				pairs: []Pair{pair1, pair2},
			},
			want: []Pair{pair1, pair2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				pairs: tt.fields.pairs,
			}
			if got := a.Pairs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchID_ID(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	type fields struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns id",
			fields: fields{
				id: id1,
			},
			want: id1,
		},
		{
			name: "returns another id",
			fields: fields{
				id: id2,
			},
			want: id2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MatchID{
				id: tt.fields.id,
			}
			if got := m.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	pair1 := makeValidPair(t)
	pair2 := makeValidPair(t)

	dupPrompt := mustPromptID(t, uuid.New())
	dupMatch1 := mustMatchID(t, uuid.New())
	dupMatch2 := mustMatchID(t, uuid.New())

	dupMatch := mustMatchID(t, uuid.New())
	dupPrompt1 := mustPromptID(t, uuid.New())
	dupPrompt2 := mustPromptID(t, uuid.New())

	type args struct {
		pairs []Pair
	}
	tests := []struct {
		name    string
		args    args
		want    Answer
		wantErr bool
	}{
		{
			name: "valid empty answer",
			args: args{
				pairs: []Pair{},
			},
			want: Answer{
				pairs: []Pair{},
			},
			wantErr: false,
		},
		{
			name: "valid answer with pairs",
			args: args{
				pairs: []Pair{pair1, pair2},
			},
			want: Answer{
				pairs: []Pair{pair1, pair2},
			},
			wantErr: false,
		},
		{
			name: "invalid duplicate prompt ids",
			args: args{
				pairs: []Pair{
					{PromptID: dupPrompt, MatchID: dupMatch1},
					{PromptID: dupPrompt, MatchID: dupMatch2},
				},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "invalid duplicate match ids",
			args: args{
				pairs: []Pair{
					{PromptID: dupPrompt1, MatchID: dupMatch},
					{PromptID: dupPrompt2, MatchID: dupMatch},
				},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "invalid nil prompt id",
			args: args{
				pairs: []Pair{
					{
						PromptID: PromptID{},
						MatchID:  mustMatchID(t, uuid.New()),
					},
				},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "invalid nil match id",
			args: args{
				pairs: []Pair{
					{
						PromptID: mustPromptID(t, uuid.New()),
						MatchID:  MatchID{},
					},
				},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "invalid too many pairs",
			args: args{
				pairs: makeValidPairs(t, matching.MaxPairs+1),
			},
			want:    Answer{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.pairs)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMatchID(t *testing.T) {
	validID := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    MatchID
		wantErr bool
	}{
		{
			name: "valid id",
			args: args{
				id: validID,
			},
			want: MatchID{
				id: validID,
			},
			wantErr: false,
		},
		{
			name: "nil id",
			args: args{
				id: uuid.Nil,
			},
			want:    MatchID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatchID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatchID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMatchID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewPromptID(t *testing.T) {
	validID := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    PromptID
		wantErr bool
	}{
		{
			name: "valid id",
			args: args{
				id: validID,
			},
			want: PromptID{
				id: validID,
			},
			wantErr: false,
		},
		{
			name: "nil id",
			args: args{
				id: uuid.Nil,
			},
			want:    PromptID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPromptID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPromptID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPromptID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromptID_ID(t *testing.T) {
	id1 := uuid.New()
	id2 := uuid.New()

	type fields struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns id",
			fields: fields{
				id: id1,
			},
			want: id1,
		},
		{
			name: "returns another id",
			fields: fields{
				id: id2,
			},
			want: id2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PromptID{
				id: tt.fields.id,
			}
			if got := p.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateAnswerPairs(t *testing.T) {
	dupPrompt := mustPromptID(t, uuid.New())
	dupMatch1 := mustMatchID(t, uuid.New())
	dupMatch2 := mustMatchID(t, uuid.New())

	type args struct {
		pairs []Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid pairs",
			args: args{
				pairs: makeValidPairs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "empty pairs valid",
			args: args{
				pairs: []Pair{},
			},
			wantErr: false,
		},
		{
			name: "invalid nil prompt",
			args: args{
				pairs: []Pair{
					{
						PromptID: PromptID{},
						MatchID:  mustMatchID(t, uuid.New()),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid duplicate prompt",
			args: args{
				pairs: []Pair{
					{PromptID: dupPrompt, MatchID: dupMatch1},
					{PromptID: dupPrompt, MatchID: dupMatch2},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid max count exceeded",
			args: args{
				pairs: makeValidPairs(t, matching.MaxPairs+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAnswerPairs(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validateAnswerPairs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateAnswerPairsMaxCount(t *testing.T) {
	type args struct {
		pairs []Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty pairs",
			args: args{
				pairs: []Pair{},
			},
			wantErr: false,
		},
		{
			name: "exact max count",
			args: args{
				pairs: makeValidPairs(t, matching.MaxPairs),
			},
			wantErr: false,
		},
		{
			name: "more than max count",
			args: args{
				pairs: makeValidPairs(t, matching.MaxPairs+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAnswerPairsMaxCount(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validateAnswerPairsMaxCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairIDsNotEmpty(t *testing.T) {
	type args struct {
		pairs []Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid ids",
			args: args{
				pairs: makeValidPairs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "empty slice",
			args: args{
				pairs: []Pair{},
			},
			wantErr: false,
		},
		{
			name: "nil prompt id",
			args: args{
				pairs: []Pair{
					{
						PromptID: PromptID{},
						MatchID:  mustMatchID(t, uuid.New()),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "nil match id",
			args: args{
				pairs: []Pair{
					{
						PromptID: mustPromptID(t, uuid.New()),
						MatchID:  MatchID{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairIDsNotEmpty(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validatePairIDsNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairsDuplicates(t *testing.T) {
	dupPrompt := mustPromptID(t, uuid.New())
	dupMatch1 := mustMatchID(t, uuid.New())
	dupMatch2 := mustMatchID(t, uuid.New())

	dupMatch := mustMatchID(t, uuid.New())
	dupPrompt1 := mustPromptID(t, uuid.New())
	dupPrompt2 := mustPromptID(t, uuid.New())

	type args struct {
		pairs []Pair
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid unique pairs",
			args: args{
				pairs: makeValidPairs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "empty slice",
			args: args{
				pairs: []Pair{},
			},
			wantErr: false,
		},
		{
			name: "duplicate prompt ids",
			args: args{
				pairs: []Pair{
					{PromptID: dupPrompt, MatchID: dupMatch1},
					{PromptID: dupPrompt, MatchID: dupMatch2},
				},
			},
			wantErr: true,
		},
		{
			name: "duplicate match ids",
			args: args{
				pairs: []Pair{
					{PromptID: dupPrompt1, MatchID: dupMatch},
					{PromptID: dupPrompt2, MatchID: dupMatch},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairsDuplicates(tt.args.pairs); (err != nil) != tt.wantErr {
				t.Errorf("validatePairsDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
