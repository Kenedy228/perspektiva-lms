package answer

import (
	"errors"
	"reflect"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	"github.com/google/uuid"
)

func mustOptionID(t *testing.T, id uuid.UUID) OptionID {
	t.Helper()

	got, err := NewOptionID(id)
	if err != nil {
		t.Fatalf("NewOptionID() error = %v", err)
	}

	return got
}

func makeOptionIDs(t *testing.T, n int) []OptionID {
	t.Helper()

	ids := make([]OptionID, 0, n)
	for i := 0; i < n; i++ {
		ids = append(ids, mustOptionID(t, uuid.New()))
	}

	return ids
}

func TestAnswer_Clone(t *testing.T) {
	optionIDs := makeOptionIDs(t, 2)

	type fields struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name   string
		fields fields
		want   question.Answer
	}{
		{
			name: "clones answer",
			fields: fields{
				optionIDs: optionIDs,
			},
			want: Answer{
				optionIDs: optionIDs,
			},
		},
		{
			name: "clones empty answer",
			fields: fields{
				optionIDs: []OptionID{},
			},
			want: Answer{
				optionIDs: []OptionID{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				optionIDs: tt.fields.optionIDs,
			}
			if got := a.Clone(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_IsEmpty(t *testing.T) {
	type fields struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty answer",
			fields: fields{
				optionIDs: []OptionID{},
			},
			want: true,
		},
		{
			name: "non-empty answer",
			fields: fields{
				optionIDs: makeOptionIDs(t, 1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				optionIDs: tt.fields.optionIDs,
			}
			if got := a.IsEmpty(); got != tt.want {
				t.Errorf("IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_OptionIDSet(t *testing.T) {
	ids := makeOptionIDs(t, 2)

	type fields struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name   string
		fields fields
		want   map[uuid.UUID]struct{}
	}{
		{
			name: "returns set of ids",
			fields: fields{
				optionIDs: ids,
			},
			want: map[uuid.UUID]struct{}{
				ids[0].ID(): {},
				ids[1].ID(): {},
			},
		},
		{
			name: "returns empty set for empty answer",
			fields: fields{
				optionIDs: []OptionID{},
			},
			want: map[uuid.UUID]struct{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				optionIDs: tt.fields.optionIDs,
			}
			if got := a.OptionIDSet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptionIDSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnswer_OptionIDs(t *testing.T) {
	ids := makeOptionIDs(t, 2)

	type fields struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name   string
		fields fields
		want   []OptionID
	}{
		{
			name: "returns option ids",
			fields: fields{
				optionIDs: ids,
			},
			want: ids,
		},
		{
			name: "returns empty slice",
			fields: fields{
				optionIDs: []OptionID{},
			},
			want: []OptionID{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Answer{
				optionIDs: tt.fields.optionIDs,
			}
			if got := a.OptionIDs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OptionIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	validIDs := makeOptionIDs(t, 2)
	dupID := mustOptionID(t, uuid.New())

	type args struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name    string
		args    args
		want    Answer
		wantErr bool
	}{
		{
			name: "creates answer with valid option ids",
			args: args{
				optionIDs: validIDs,
			},
			want: Answer{
				optionIDs: validIDs,
			},
			wantErr: false,
		},
		{
			name: "creates empty answer",
			args: args{
				optionIDs: []OptionID{},
			},
			want: Answer{
				optionIDs: []OptionID{},
			},
			wantErr: false,
		},
		{
			name: "returns error for empty option id",
			args: args{
				optionIDs: []OptionID{{}},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "returns error for duplicate option ids",
			args: args{
				optionIDs: []OptionID{dupID, dupID},
			},
			want:    Answer{},
			wantErr: true,
		},
		{
			name: "returns error for too many option ids",
			args: args{
				optionIDs: makeOptionIDs(t, selectable.MaxOptionsCount+1),
			},
			want:    Answer{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.optionIDs)
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

func TestNewOptionID(t *testing.T) {
	validID := uuid.New()

	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		want    OptionID
		wantErr bool
	}{
		{
			name: "creates option id",
			args: args{
				id: validID,
			},
			want:    OptionID{id: validID},
			wantErr: false,
		},
		{
			name: "returns error for nil uuid",
			args: args{
				id: uuid.Nil,
			},
			want:    OptionID{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewOptionID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewOptionID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("NewOptionID() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewOptionID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionID_ID(t *testing.T) {
	id := uuid.New()

	type fields struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns uuid",
			fields: fields{
				id: id,
			},
			want: id,
		},
		{
			name: "returns nil uuid for zero value",
			fields: fields{
				id: uuid.Nil,
			},
			want: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := OptionID{
				id: tt.fields.id,
			}
			if got := o.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptionID_IsZero(t *testing.T) {
	type fields struct {
		id uuid.UUID
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "false for valid id",
			fields: fields{
				id: uuid.New(),
			},
			want: false,
		},
		{
			name: "true for nil id",
			fields: fields{
				id: uuid.Nil,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := OptionID{
				id: tt.fields.id,
			}
			if got := o.IsZero(); got != tt.want {
				t.Errorf("IsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateOptionIDs(t *testing.T) {
	dupID := mustOptionID(t, uuid.New())

	type args struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid option ids",
			args: args{
				optionIDs: makeOptionIDs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "empty slice valid",
			args: args{
				optionIDs: []OptionID{},
			},
			wantErr: false,
		},
		{
			name: "contains empty option id",
			args: args{
				optionIDs: []OptionID{{}},
			},
			wantErr: true,
		},
		{
			name: "contains duplicates",
			args: args{
				optionIDs: []OptionID{dupID, dupID},
			},
			wantErr: true,
		},
		{
			name: "too many option ids",
			args: args{
				optionIDs: makeOptionIDs(t, selectable.MaxOptionsCount+1),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionIDs(tt.args.optionIDs); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionIDs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionIDsContainsEmpty(t *testing.T) {
	type args struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "all option ids valid",
			args: args{
				optionIDs: makeOptionIDs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "contains zero option id",
			args: args{
				optionIDs: []OptionID{{}},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				optionIDs: []OptionID{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionIDsContainsEmpty(tt.args.optionIDs); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionIDsContainsEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionIDsCount(t *testing.T) {
	type args struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "within limit",
			args: args{
				optionIDs: makeOptionIDs(t, selectable.MaxOptionsCount),
			},
			wantErr: false,
		},
		{
			name: "exceeds limit",
			args: args{
				optionIDs: makeOptionIDs(t, selectable.MaxOptionsCount+1),
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				optionIDs: []OptionID{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionIDsCount(tt.args.optionIDs); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionIDsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateOptionIDsDuplicates(t *testing.T) {
	dupID := mustOptionID(t, uuid.New())

	type args struct {
		optionIDs []OptionID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "unique option ids",
			args: args{
				optionIDs: makeOptionIDs(t, 2),
			},
			wantErr: false,
		},
		{
			name: "duplicate option ids",
			args: args{
				optionIDs: []OptionID{dupID, dupID},
			},
			wantErr: true,
		},
		{
			name: "empty slice",
			args: args{
				optionIDs: []OptionID{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateOptionIDsDuplicates(tt.args.optionIDs); (err != nil) != tt.wantErr {
				t.Errorf("validateOptionIDsDuplicates() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
