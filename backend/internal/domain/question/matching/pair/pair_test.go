package pair

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/shared/text"
	"github.com/google/uuid"
)

func mustText(t *testing.T, value string) text.Text {
	t.Helper()

	got, err := text.New(value)
	if err != nil {
		t.Fatalf("text.New() error = %v", err)
	}

	return got
}

func mustPrompt(t *testing.T, value string) Prompt {
	t.Helper()

	got, err := NewPrompt(mustText(t, value))
	if err != nil {
		t.Fatalf("NewPrompt() error = %v", err)
	}

	return got
}

func mustMatch(t *testing.T, value string) Match {
	t.Helper()

	got, err := NewMatch(mustText(t, value))
	if err != nil {
		t.Fatalf("NewMatch() error = %v", err)
	}

	return got
}

func mustPair(t *testing.T, promptValue, matchValue string) Pair {
	t.Helper()

	got, err := New(
		mustPrompt(t, promptValue),
		mustMatch(t, matchValue),
	)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return got
}

func TestMatch_ID(t *testing.T) {
	text1 := mustText(t, "match")
	id1 := uuid.New()

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns match id",
			fields: fields{
				id:   id1,
				text: text1,
			},
			want: id1,
		},
		{
			name: "returns nil for zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Match{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := m.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatch_IsIncomplete(t *testing.T) {
	validText := mustText(t, "valid match")
	validID := uuid.New()

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "complete match",
			fields: fields{
				id:   validID,
				text: validText,
			},
			want: false,
		},
		{
			name: "incomplete because nil id",
			fields: fields{
				id:   uuid.Nil,
				text: validText,
			},
			want: true,
		},
		{
			name: "incomplete because empty text",
			fields: fields{
				id:   validID,
				text: text.Text{},
			},
			want: true,
		},
		{
			name: "incomplete zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Match{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := m.IsIncomplete(); got != tt.want {
				t.Errorf("IsIncomplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatch_Text(t *testing.T) {
	text1 := mustText(t, "match")

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   text.Text
	}{
		{
			name: "returns text",
			fields: fields{
				id:   uuid.New(),
				text: text1,
			},
			want: text1,
		},
		{
			name: "returns zero text for zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: text.Text{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Match{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := m.Text(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	validPrompt := mustPrompt(t, "prompt")
	validMatch := mustMatch(t, "match")

	type args struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name    string
		args    args
		want    Pair
		wantErr bool
	}{
		{
			name: "valid pair",
			args: args{
				prompt: validPrompt,
				match:  validMatch,
			},
			want: Pair{
				prompt: validPrompt,
				match:  validMatch,
			},
			wantErr: false,
		},
		{
			name: "invalid empty prompt",
			args: args{
				prompt: Prompt{},
				match:  validMatch,
			},
			want:    Pair{},
			wantErr: true,
		},
		{
			name: "invalid empty match",
			args: args{
				prompt: validPrompt,
				match:  Match{},
			},
			want:    Pair{},
			wantErr: true,
		},
		{
			name: "invalid both empty",
			args: args{
				prompt: Prompt{},
				match:  Match{},
			},
			want:    Pair{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.prompt, tt.args.match)
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

func TestNewMatch(t *testing.T) {
	validText := mustText(t, "match")
	longText := mustText(t, strings.Repeat("а", MatchCharsLimit+1))

	type args struct {
		t text.Text
	}
	tests := []struct {
		name     string
		args     args
		wantText text.Text
		wantErr  bool
	}{
		{
			name: "valid match",
			args: args{
				t: validText,
			},
			wantText: validText,
			wantErr:  false,
		},
		{
			name: "too long match",
			args: args{
				t: longText,
			},
			wantText: text.Text{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatch(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("NewMatch() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}

			if got.ID() == uuid.Nil {
				t.Errorf("NewMatch() ID() = %v, want non-nil uuid", got.ID())
			}
			if !reflect.DeepEqual(got.Text(), tt.wantText) {
				t.Errorf("NewMatch() Text() = %v, want %v", got.Text(), tt.wantText)
			}
		})
	}
}

func TestNewPrompt(t *testing.T) {
	validText := mustText(t, "prompt")
	longText := mustText(t, strings.Repeat("а", PromptCharsLimit+1))

	type args struct {
		t text.Text
	}
	tests := []struct {
		name     string
		args     args
		wantText text.Text
		wantErr  bool
	}{
		{
			name: "valid prompt",
			args: args{
				t: validText,
			},
			wantText: validText,
			wantErr:  false,
		},
		{
			name: "too long prompt",
			args: args{
				t: longText,
			},
			wantText: text.Text{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrompt(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Errorf("NewPrompt() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}

			if got.ID() == uuid.Nil {
				t.Errorf("NewPrompt() ID() = %v, want non-nil uuid", got.ID())
			}
			if !reflect.DeepEqual(got.Text(), tt.wantText) {
				t.Errorf("NewPrompt() Text() = %v, want %v", got.Text(), tt.wantText)
			}
		})
	}
}

func TestPair_IsIncomplete(t *testing.T) {
	validPrompt := mustPrompt(t, "prompt")
	validMatch := mustMatch(t, "match")

	type fields struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "complete pair",
			fields: fields{
				prompt: validPrompt,
				match:  validMatch,
			},
			want: false,
		},
		{
			name: "incomplete because prompt empty",
			fields: fields{
				prompt: Prompt{},
				match:  validMatch,
			},
			want: true,
		},
		{
			name: "incomplete because match empty",
			fields: fields{
				prompt: validPrompt,
				match:  Match{},
			},
			want: true,
		},
		{
			name: "incomplete because both empty",
			fields: fields{
				prompt: Prompt{},
				match:  Match{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pair{
				prompt: tt.fields.prompt,
				match:  tt.fields.match,
			}
			if got := p.IsIncomplete(); got != tt.want {
				t.Errorf("IsIncomplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPair_Match(t *testing.T) {
	pair1 := mustPair(t, "prompt", "match")

	type fields struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name   string
		fields fields
		want   Match
	}{
		{
			name: "returns match",
			fields: fields{
				prompt: pair1.Prompt(),
				match:  pair1.Match(),
			},
			want: pair1.Match(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pair{
				prompt: tt.fields.prompt,
				match:  tt.fields.match,
			}
			if got := p.Match(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPair_MatchID(t *testing.T) {
	pair1 := mustPair(t, "prompt", "match")

	type fields struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns match id",
			fields: fields{
				prompt: pair1.Prompt(),
				match:  pair1.Match(),
			},
			want: pair1.Match().ID(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pair{
				prompt: tt.fields.prompt,
				match:  tt.fields.match,
			}
			if got := p.MatchID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MatchID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPair_Prompt(t *testing.T) {
	pair1 := mustPair(t, "prompt", "match")

	type fields struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name   string
		fields fields
		want   Prompt
	}{
		{
			name: "returns prompt",
			fields: fields{
				prompt: pair1.Prompt(),
				match:  pair1.Match(),
			},
			want: pair1.Prompt(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pair{
				prompt: tt.fields.prompt,
				match:  tt.fields.match,
			}
			if got := p.Prompt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPair_PromptID(t *testing.T) {
	pair1 := mustPair(t, "prompt", "match")

	type fields struct {
		prompt Prompt
		match  Match
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns prompt id",
			fields: fields{
				prompt: pair1.Prompt(),
				match:  pair1.Match(),
			},
			want: pair1.Prompt().ID(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pair{
				prompt: tt.fields.prompt,
				match:  tt.fields.match,
			}
			if got := p.PromptID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PromptID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_ID(t *testing.T) {
	text1 := mustText(t, "prompt")
	id1 := uuid.New()

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   uuid.UUID
	}{
		{
			name: "returns prompt id",
			fields: fields{
				id:   id1,
				text: text1,
			},
			want: id1,
		},
		{
			name: "returns nil for zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: uuid.Nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prompt{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := p.ID(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_IsIncomplete(t *testing.T) {
	validText := mustText(t, "valid prompt")
	validID := uuid.New()

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "complete prompt",
			fields: fields{
				id:   validID,
				text: validText,
			},
			want: false,
		},
		{
			name: "incomplete because nil id",
			fields: fields{
				id:   uuid.Nil,
				text: validText,
			},
			want: true,
		},
		{
			name: "incomplete because empty text",
			fields: fields{
				id:   validID,
				text: text.Text{},
			},
			want: true,
		},
		{
			name: "incomplete zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prompt{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := p.IsIncomplete(); got != tt.want {
				t.Errorf("IsIncomplete() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrompt_Text(t *testing.T) {
	text1 := mustText(t, "prompt")

	type fields struct {
		id   uuid.UUID
		text text.Text
	}
	tests := []struct {
		name   string
		fields fields
		want   text.Text
	}{
		{
			name: "returns text",
			fields: fields{
				id:   uuid.New(),
				text: text1,
			},
			want: text1,
		},
		{
			name: "returns zero text for zero value",
			fields: fields{
				id:   uuid.Nil,
				text: text.Text{},
			},
			want: text.Text{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Prompt{
				id:   tt.fields.id,
				text: tt.fields.text,
			}
			if got := p.Text(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Text() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateMatch(t *testing.T) {
	validMatch := mustMatch(t, "match")

	type args struct {
		match Match
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid match",
			args: args{
				match: validMatch,
			},
			wantErr: false,
		},
		{
			name: "zero match",
			args: args{
				match: Match{},
			},
			wantErr: true,
		},
		{
			name: "match with nil id",
			args: args{
				match: Match{
					id:   uuid.Nil,
					text: mustText(t, "match"),
				},
			},
			wantErr: true,
		},
		{
			name: "match with empty text",
			args: args{
				match: Match{
					id:   uuid.New(),
					text: text.Text{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateMatch(tt.args.match); (err != nil) != tt.wantErr {
				t.Errorf("validateMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePrompt(t *testing.T) {
	validPrompt := mustPrompt(t, "prompt")

	type args struct {
		prompt Prompt
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid prompt",
			args: args{
				prompt: validPrompt,
			},
			wantErr: false,
		},
		{
			name: "zero prompt",
			args: args{
				prompt: Prompt{},
			},
			wantErr: true,
		},
		{
			name: "prompt with nil id",
			args: args{
				prompt: Prompt{
					id:   uuid.Nil,
					text: mustText(t, "prompt"),
				},
			},
			wantErr: true,
		},
		{
			name: "prompt with empty text",
			args: args{
				prompt: Prompt{
					id:   uuid.New(),
					text: text.Text{},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validatePrompt(tt.args.prompt); (err != nil) != tt.wantErr {
				t.Errorf("validatePrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
