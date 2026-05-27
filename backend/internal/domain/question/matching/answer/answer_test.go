package answer

import (
	"errors"
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/matching"
	"github.com/google/uuid"
)

func makeValidPair() Pair {
	return Pair{
		PromptID: uuid.New(),
		MatchID:  uuid.New(),
	}
}

func makeValidPairs(n int) []Pair {
	pairs := make([]Pair, 0, n)
	for i := 0; i < n; i++ {
		pairs = append(pairs, makeValidPair())
	}
	return pairs
}

func TestNew(t *testing.T) {
	pair1 := makeValidPair()
	pair2 := makeValidPair()

	dupPrompt := uuid.New()
	dupMatchLeft := uuid.New()
	dupMatchRight := uuid.New()

	dupMatch := uuid.New()
	dupPromptLeft := uuid.New()
	dupPromptRight := uuid.New()

	tests := []struct {
		name    string
		pairs   []Pair
		want    Answer
		wantErr bool
	}{
		{
			name:    "valid empty answer",
			pairs:   []Pair{},
			want:    Answer{pairs: []Pair{}},
			wantErr: false,
		},
		{
			name:    "valid answer with pairs",
			pairs:   []Pair{pair1, pair2},
			want:    Answer{pairs: []Pair{pair1, pair2}},
			wantErr: false,
		},
		{
			name: "invalid duplicate prompt id",
			pairs: []Pair{
				{PromptID: dupPrompt, MatchID: dupMatchLeft},
				{PromptID: dupPrompt, MatchID: dupMatchRight},
			},
			wantErr: true,
		},
		{
			name: "invalid duplicate match id",
			pairs: []Pair{
				{PromptID: dupPromptLeft, MatchID: dupMatch},
				{PromptID: dupPromptRight, MatchID: dupMatch},
			},
			wantErr: true,
		},
		{
			name: "invalid nil prompt id",
			pairs: []Pair{
				{PromptID: uuid.Nil, MatchID: uuid.New()},
			},
			wantErr: true,
		},
		{
			name: "invalid nil match id",
			pairs: []Pair{
				{PromptID: uuid.New(), MatchID: uuid.Nil},
			},
			wantErr: true,
		},
		{
			name:    "invalid too many pairs",
			pairs:   makeValidPairs(matching.MaxPairs + 1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.pairs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("New() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_Pairs(t *testing.T) {
	tests := []struct {
		name  string
		pairs []Pair
	}{
		{name: "empty", pairs: []Pair{}},
		{name: "non-empty", pairs: makeValidPairs(2)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{pairs: tt.pairs}
			got := a.Pairs()

			if !reflect.DeepEqual(got, tt.pairs) {
				t.Fatalf("Pairs() = %v, want %v", got, tt.pairs)
			}

			if len(got) > 0 {
				got[0] = Pair{PromptID: uuid.New(), MatchID: uuid.New()}
				if reflect.DeepEqual(got, a.pairs) {
					t.Fatalf("Pairs() must return a cloned slice")
				}
			}
		})
	}
}

func TestAnswer_AsMap(t *testing.T) {
	pair1 := makeValidPair()
	pair2 := makeValidPair()

	tests := []struct {
		name  string
		pairs []Pair
		want  map[uuid.UUID]uuid.UUID
	}{
		{
			name:  "empty pairs",
			pairs: []Pair{},
			want:  map[uuid.UUID]uuid.UUID{},
		},
		{
			name:  "non-empty pairs",
			pairs: []Pair{pair1, pair2},
			want: map[uuid.UUID]uuid.UUID{
				pair1.PromptID: pair1.MatchID,
				pair2.PromptID: pair2.MatchID,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{pairs: tt.pairs}
			got := a.AsMap()
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("AsMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_IsEmpty(t *testing.T) {
	tests := []struct {
		name  string
		pairs []Pair
		want  bool
	}{
		{name: "empty", pairs: []Pair{}, want: true},
		{name: "non-empty", pairs: []Pair{makeValidPair()}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{pairs: tt.pairs}
			if got := a.IsEmpty(); got != tt.want {
				t.Fatalf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_Clone(t *testing.T) {
	pairs := makeValidPairs(2)
	a := Answer{pairs: pairs}

	got, ok := a.Clone().(Answer)
	if !ok {
		t.Fatalf("Clone() type = %T, want Answer", a.Clone())
	}
	if !reflect.DeepEqual(got, a) {
		t.Fatalf("Clone() = %v, want %v", got, a)
	}
	if len(got.pairs) > 0 {
		got.pairs[0] = Pair{PromptID: uuid.New(), MatchID: uuid.New()}
		if reflect.DeepEqual(got, a) {
			t.Fatalf("Clone() must return an independent copy")
		}
	}

	var _ question.Answer = got
}

func Test_validateAnswerPairs(t *testing.T) {
	dupPrompt := uuid.New()
	dupMatchLeft := uuid.New()
	dupMatchRight := uuid.New()

	tests := []struct {
		name    string
		pairs   []Pair
		wantErr bool
	}{
		{name: "valid pairs", pairs: makeValidPairs(2), wantErr: false},
		{name: "empty pairs", pairs: []Pair{}, wantErr: false},
		{
			name: "nil prompt id",
			pairs: []Pair{
				{PromptID: uuid.Nil, MatchID: uuid.New()},
			},
			wantErr: true,
		},
		{
			name: "duplicate prompt id",
			pairs: []Pair{
				{PromptID: dupPrompt, MatchID: dupMatchLeft},
				{PromptID: dupPrompt, MatchID: dupMatchRight},
			},
			wantErr: true,
		},
		{
			name:    "max count exceeded",
			pairs:   makeValidPairs(matching.MaxPairs + 1),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAnswerPairs(tt.pairs); (err != nil) != tt.wantErr {
				t.Fatalf("validateAnswerPairs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateAnswerPairsMaxCount(t *testing.T) {
	tests := []struct {
		name    string
		pairs   []Pair
		wantErr bool
	}{
		{name: "empty", pairs: []Pair{}, wantErr: false},
		{name: "exact max", pairs: makeValidPairs(matching.MaxPairs), wantErr: false},
		{name: "above max", pairs: makeValidPairs(matching.MaxPairs + 1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateAnswerPairsMaxCount(tt.pairs); (err != nil) != tt.wantErr {
				t.Fatalf("validateAnswerPairsMaxCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairIDsNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		pairs   []Pair
		wantErr bool
	}{
		{name: "valid ids", pairs: makeValidPairs(2), wantErr: false},
		{name: "empty slice", pairs: []Pair{}, wantErr: false},
		{
			name: "nil prompt id",
			pairs: []Pair{
				{PromptID: uuid.Nil, MatchID: uuid.New()},
			},
			wantErr: true,
		},
		{
			name: "nil match id",
			pairs: []Pair{
				{PromptID: uuid.New(), MatchID: uuid.Nil},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairIDsNotEmpty(tt.pairs); (err != nil) != tt.wantErr {
				t.Fatalf("validatePairIDsNotEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePairsDuplicates(t *testing.T) {
	dupPrompt := uuid.New()
	dupMatchLeft := uuid.New()
	dupMatchRight := uuid.New()

	dupMatch := uuid.New()
	dupPromptLeft := uuid.New()
	dupPromptRight := uuid.New()

	tests := []struct {
		name    string
		pairs   []Pair
		wantErr bool
	}{
		{name: "valid unique pairs", pairs: makeValidPairs(2), wantErr: false},
		{name: "empty slice", pairs: []Pair{}, wantErr: false},
		{
			name: "duplicate prompt ids",
			pairs: []Pair{
				{PromptID: dupPrompt, MatchID: dupMatchLeft},
				{PromptID: dupPrompt, MatchID: dupMatchRight},
			},
			wantErr: true,
		},
		{
			name: "duplicate match ids",
			pairs: []Pair{
				{PromptID: dupPromptLeft, MatchID: dupMatch},
				{PromptID: dupPromptRight, MatchID: dupMatch},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePairsDuplicates(tt.pairs); (err != nil) != tt.wantErr {
				t.Fatalf("validatePairsDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
