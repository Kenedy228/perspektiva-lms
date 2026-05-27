package pair

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func mustPrompt(t *testing.T, value string) Prompt {
	t.Helper()

	got, err := NewPrompt(value)
	if err != nil {
		t.Fatalf("NewPrompt() error = %v", err)
	}

	return got
}

func mustMatch(t *testing.T, value string) Match {
	t.Helper()

	got, err := NewMatch(value)
	if err != nil {
		t.Fatalf("NewMatch() error = %v", err)
	}

	return got
}

func mustPair(t *testing.T, promptValue, matchValue string) Pair {
	t.Helper()

	got, err := New(mustPrompt(t, promptValue), mustMatch(t, matchValue))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	return got
}

func TestNewPrompt(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid prompt", value: "prompt"},
		{name: "valid unicode prompt", value: "Привет, мир"},
		{name: "value at limit", value: strings.Repeat("я", ValueCharsLimit)},
		{name: "empty value", value: "", wantErr: true},
		{name: "value over limit", value: strings.Repeat("я", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPrompt(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewPrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("NewPrompt() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if got.ID() == uuid.Nil {
				t.Fatalf("NewPrompt() ID() = nil, want non-nil uuid")
			}
			if got.Value() != tt.value {
				t.Fatalf("NewPrompt() Value() = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}

func TestRestorePrompt(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		value   string
		want    Prompt
		wantErr bool
	}{
		{name: "valid restore", id: id, value: "prompt", want: Prompt{id: id, value: "prompt"}},
		{name: "nil id", id: uuid.Nil, value: "prompt", wantErr: true},
		{name: "empty value", id: id, value: "", wantErr: true},
		{name: "over limit", id: id, value: strings.Repeat("я", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RestorePrompt(tt.id, tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RestorePrompt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("RestorePrompt() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("RestorePrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPromptMethods(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name   string
		prompt Prompt
		wantID uuid.UUID
		wantV  string
		wantZ  bool
	}{
		{name: "valid prompt", prompt: Prompt{id: id, value: "prompt"}, wantID: id, wantV: "prompt", wantZ: false},
		{name: "zero value", prompt: Prompt{}, wantID: uuid.Nil, wantV: "", wantZ: true},
		{name: "nil id", prompt: Prompt{id: uuid.Nil, value: "prompt"}, wantID: uuid.Nil, wantV: "prompt", wantZ: true},
		{name: "empty value", prompt: Prompt{id: id, value: ""}, wantID: id, wantV: "", wantZ: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.prompt.ID(); got != tt.wantID {
				t.Fatalf("ID() = %v, want %v", got, tt.wantID)
			}
			if got := tt.prompt.Value(); got != tt.wantV {
				t.Fatalf("Value() = %q, want %q", got, tt.wantV)
			}
			if got := tt.prompt.IsZero(); got != tt.wantZ {
				t.Fatalf("IsZero() = %v, want %v", got, tt.wantZ)
			}
		})
	}
}

func TestNewMatch(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "valid match", value: "match"},
		{name: "valid unicode match", value: "Совпадение"},
		{name: "value at limit", value: strings.Repeat("ы", ValueCharsLimit)},
		{name: "empty value", value: "", wantErr: true},
		{name: "value over limit", value: strings.Repeat("ы", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMatch(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("NewMatch() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if got.ID() == uuid.Nil {
				t.Fatalf("NewMatch() ID() = nil, want non-nil uuid")
			}
			if got.Value() != tt.value {
				t.Fatalf("NewMatch() Value() = %q, want %q", got.Value(), tt.value)
			}
		})
	}
}

func TestRestoreMatch(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name    string
		id      uuid.UUID
		value   string
		want    Match
		wantErr bool
	}{
		{name: "valid restore", id: id, value: "match", want: Match{id: id, value: "match"}},
		{name: "nil id", id: uuid.Nil, value: "match", wantErr: true},
		{name: "empty value", id: id, value: "", wantErr: true},
		{name: "over limit", id: id, value: strings.Repeat("ы", ValueCharsLimit+1), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RestoreMatch(tt.id, tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("RestoreMatch() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("RestoreMatch() error = %v, want wrapped ErrInvalid", err)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("RestoreMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatchMethods(t *testing.T) {
	id := uuid.New()

	tests := []struct {
		name  string
		match Match
		wantI uuid.UUID
		wantV string
		wantZ bool
	}{
		{name: "valid match", match: Match{id: id, value: "match"}, wantI: id, wantV: "match", wantZ: false},
		{name: "zero value", match: Match{}, wantI: uuid.Nil, wantV: "", wantZ: true},
		{name: "nil id", match: Match{id: uuid.Nil, value: "match"}, wantI: uuid.Nil, wantV: "match", wantZ: true},
		{name: "empty value", match: Match{id: id, value: ""}, wantI: id, wantV: "", wantZ: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.match.ID(); got != tt.wantI {
				t.Fatalf("ID() = %v, want %v", got, tt.wantI)
			}
			if got := tt.match.Value(); got != tt.wantV {
				t.Fatalf("Value() = %q, want %q", got, tt.wantV)
			}
			if got := tt.match.IsZero(); got != tt.wantZ {
				t.Fatalf("IsZero() = %v, want %v", got, tt.wantZ)
			}
		})
	}
}

func TestNewPair(t *testing.T) {
	validPrompt := mustPrompt(t, "prompt")
	validMatch := mustMatch(t, "match")

	tests := []struct {
		name    string
		prompt  Prompt
		match   Match
		want    Pair
		wantErr bool
	}{
		{
			name:    "valid pair",
			prompt:  validPrompt,
			match:   validMatch,
			want:    Pair{prompt: validPrompt, match: validMatch},
			wantErr: false,
		},
		{name: "zero prompt", prompt: Prompt{}, match: validMatch, wantErr: true},
		{name: "zero match", prompt: validPrompt, match: Match{}, wantErr: true},
		{name: "both zero", prompt: Prompt{}, match: Match{}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.prompt, tt.match)
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

func TestPairMethods(t *testing.T) {
	pairValue := mustPair(t, "prompt", "match")

	tests := []struct {
		name string
		pair Pair
		zero bool
	}{
		{name: "valid pair", pair: pairValue, zero: false},
		{name: "zero prompt", pair: Pair{prompt: Prompt{}, match: pairValue.Match()}, zero: true},
		{name: "zero match", pair: Pair{prompt: pairValue.Prompt(), match: Match{}}, zero: true},
		{name: "zero pair", pair: Pair{}, zero: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pair.IsZero(); got != tt.zero {
				t.Fatalf("IsZero() = %v, want %v", got, tt.zero)
			}

			if !tt.zero {
				if got := tt.pair.PromptID(); got != tt.pair.Prompt().ID() {
					t.Fatalf("PromptID() = %v, want %v", got, tt.pair.Prompt().ID())
				}
				if got := tt.pair.MatchID(); got != tt.pair.Match().ID() {
					t.Fatalf("MatchID() = %v, want %v", got, tt.pair.Match().ID())
				}
			}
		})
	}
}

func TestValidatorHelpers(t *testing.T) {
	tests := []struct {
		name string
		run  func() error
		want bool
	}{
		{name: "validateRequired valid", run: func() error { return validateRequired("abc") }, want: false},
		{name: "validateRequired empty", run: func() error { return validateRequired("") }, want: true},
		{name: "validateCharsLimit valid unicode", run: func() error { return validateCharsLimit(strings.Repeat("ж", ValueCharsLimit)) }, want: false},
		{name: "validateCharsLimit over unicode", run: func() error { return validateCharsLimit(strings.Repeat("ж", ValueCharsLimit+1)) }, want: true},
		{name: "validateIDRequired valid", run: func() error { return validateIDRequired(uuid.New()) }, want: false},
		{name: "validateIDRequired nil", run: func() error { return validateIDRequired(uuid.Nil) }, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if (err != nil) != tt.want {
				t.Fatalf("validator error = %v, wantErr %v", err, tt.want)
			}
			if tt.want && !errors.Is(err, ErrInvalid) {
				t.Fatalf("validator error = %v, want wrapped ErrInvalid", err)
			}
		})
	}
}

func Test_validatePromptAndMatch(t *testing.T) {
	validPrompt := mustPrompt(t, "prompt")
	validMatch := mustMatch(t, "match")

	tests := []struct {
		name    string
		run     func() error
		wantErr bool
	}{
		{name: "validatePrompt valid", run: func() error { return validatePrompt(validPrompt) }, wantErr: false},
		{name: "validatePrompt zero", run: func() error { return validatePrompt(Prompt{}) }, wantErr: true},
		{name: "validateMatch valid", run: func() error { return validateMatch(validMatch) }, wantErr: false},
		{name: "validateMatch zero", run: func() error { return validateMatch(Match{}) }, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.run()
			if (err != nil) != tt.wantErr {
				t.Fatalf("validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
