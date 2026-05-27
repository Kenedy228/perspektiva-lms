package answer

import (
	"errors"
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	"github.com/google/uuid"
)

func makeOptionIDs(n int) []uuid.UUID {
	ids := make([]uuid.UUID, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, uuid.New())
	}
	return ids
}

func TestNew(t *testing.T) {
	validIDs := makeOptionIDs(2)
	dupID := uuid.New()

	tests := []struct {
		name      string
		optionIDs []uuid.UUID
		want      Answer
		wantErr   bool
	}{
		{
			name:      "creates answer with selected options",
			optionIDs: validIDs,
			want:      Answer{optionIDs: validIDs},
		},
		{
			name:      "creates empty answer",
			optionIDs: []uuid.UUID{},
			want:      Answer{optionIDs: []uuid.UUID{}},
		},
		{
			name:      "returns error for empty option id",
			optionIDs: []uuid.UUID{uuid.Nil},
			wantErr:   true,
		},
		{
			name:      "returns error for duplicate option ids",
			optionIDs: []uuid.UUID{dupID, dupID},
			wantErr:   true,
		},
		{
			name:      "returns error for too many option ids",
			optionIDs: makeOptionIDs(selectable.MaxOptionsCount + 1),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.optionIDs)
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

func TestAnswer_OptionIDs(t *testing.T) {
	ids := makeOptionIDs(2)
	a := Answer{optionIDs: ids}

	got := a.OptionIDs()
	if !reflect.DeepEqual(got, ids) {
		t.Fatalf("OptionIDs() = %v, want %v", got, ids)
	}
	if len(got) > 0 {
		got[0] = uuid.New()
	}
	if reflect.DeepEqual(got, a.optionIDs) {
		t.Fatalf("OptionIDs() must return cloned slice")
	}
}

func TestAnswer_OptionIDSet(t *testing.T) {
	ids := makeOptionIDs(2)
	a := Answer{optionIDs: ids}

	got := a.OptionIDSet()
	want := map[uuid.UUID]struct{}{
		ids[0]: {},
		ids[1]: {},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("OptionIDSet() = %v, want %v", got, want)
	}
}

func TestAnswer_IsEmpty(t *testing.T) {
	if !((Answer{}).IsEmpty()) {
		t.Fatalf("zero answer must be empty")
	}

	if (Answer{optionIDs: makeOptionIDs(1)}).IsEmpty() {
		t.Fatalf("answer with selected option must not be empty")
	}
}

func TestAnswer_Clone(t *testing.T) {
	ids := makeOptionIDs(2)
	a := Answer{optionIDs: ids}

	got, ok := a.Clone().(Answer)
	if !ok {
		t.Fatalf("Clone() type = %T, want Answer", a.Clone())
	}
	if !reflect.DeepEqual(got, a) {
		t.Fatalf("Clone() = %v, want %v", got, a)
	}

	var _ question.Answer = got
}

func Test_validateOptionIDs(t *testing.T) {
	dupID := uuid.New()

	tests := []struct {
		name      string
		optionIDs []uuid.UUID
		wantErr   bool
	}{
		{name: "valid option ids", optionIDs: makeOptionIDs(2)},
		{name: "empty slice valid", optionIDs: []uuid.UUID{}},
		{name: "contains empty option id", optionIDs: []uuid.UUID{uuid.Nil}, wantErr: true},
		{name: "contains duplicates", optionIDs: []uuid.UUID{dupID, dupID}, wantErr: true},
		{name: "too many option ids", optionIDs: makeOptionIDs(selectable.MaxOptionsCount + 1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionIDs(tt.optionIDs); (err != nil) != tt.wantErr {
				t.Fatalf("validateOptionIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
