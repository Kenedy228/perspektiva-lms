package matching

import (
	"errors"
	"strings"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/base"
	basetitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/question/matching/pair"
)

func mustBase(t *testing.T, value string) *base.Base {
	t.Helper()

	ttl, err := basetitle.New(value)
	if err != nil {
		t.Fatalf("title.New() error = %v", err)
	}

	b, err := base.New(ttl)
	if err != nil {
		t.Fatalf("base.New() error = %v", err)
	}

	return b
}

func mustPair(t *testing.T, promptValue, matchValue string) pair.Pair {
	t.Helper()

	prompt, err := pair.NewPrompt(promptValue)
	if err != nil {
		t.Fatalf("pair.NewPrompt() error = %v", err)
	}

	match, err := pair.NewMatch(matchValue)
	if err != nil {
		t.Fatalf("pair.NewMatch() error = %v", err)
	}

	p, err := pair.New(prompt, match)
	if err != nil {
		t.Fatalf("pair.New() error = %v", err)
	}

	return p
}

func makePairs(t *testing.T, n int) []pair.Pair {
	t.Helper()

	result := make([]pair.Pair, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, mustPair(t, "prompt", "match"))
	}

	return result
}

func TestNew(t *testing.T) {
	validBase := mustBase(t, "matching question")
	validPairs := makePairs(t, MinPairs)

	tests := []struct {
		name       string
		base       *base.Base
		pairs      []pair.Pair
		wantErr    bool
		errText    string
		checkError bool
	}{
		{
			name:    "успех",
			base:    validBase,
			pairs:   validPairs,
			wantErr: false,
		},
		{
			name:       "ошибка когда base nil",
			base:       nil,
			pairs:      validPairs,
			wantErr:    true,
			errText:    "база вопросов обязательна",
			checkError: true,
		},
		{
			name:       "ошибка когда пар меньше минимума",
			base:       validBase,
			pairs:      makePairs(t, MinPairs-1),
			wantErr:    true,
			errText:    "не меньше",
			checkError: true,
		},
		{
			name: "ошибка когда есть пустая пара",
			base: validBase,
			pairs: []pair.Pair{
				mustPair(t, "p1", "m1"),
				{},
			},
			wantErr:    true,
			errText:    "под индексом 1",
			checkError: true,
		},
		{
			name:       "ошибка когда пар больше максимума",
			base:       validBase,
			pairs:      makePairs(t, MaxPairs+1),
			wantErr:    true,
			errText:    "не больше",
			checkError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.base, tt.pairs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("New() error = %v, want ErrInvalid", err)
				}
				if tt.checkError && !strings.Contains(err.Error(), tt.errText) {
					t.Fatalf("New() error text = %q, want contain %q", err.Error(), tt.errText)
				}
				return
			}

			if got == nil {
				t.Fatal("New() got nil question")
			}
			if got.Type() != question.TypeMatching {
				t.Fatalf("New() Type() = %v, want %v", got.Type(), question.TypeMatching)
			}
			if got.Instruction() != question.TypeMatching.DefaultInstruction() {
				t.Fatalf("New() Instruction() = %q, want %q", got.Instruction(), question.TypeMatching.DefaultInstruction())
			}
		})
	}
}

func TestRestore(t *testing.T) {
	validBase := mustBase(t, "matching question")
	validPairs := makePairs(t, MinPairs)

	tests := []struct {
		name    string
		base    *base.Base
		pairs   []pair.Pair
		wantErr bool
	}{
		{
			name:    "успех",
			base:    validBase,
			pairs:   validPairs,
			wantErr: false,
		},
		{
			name:    "ошибка при base nil",
			base:    nil,
			pairs:   validPairs,
			wantErr: true,
		},
		{
			name:    "ошибка при пустой паре",
			base:    validBase,
			pairs:   []pair.Pair{mustPair(t, "p1", "m1"), {}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Restore(tt.base, tt.pairs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("Restore() error = %v, want ErrInvalid", err)
				}
				return
			}
			if got == nil {
				t.Fatal("Restore() got nil question")
			}
		})
	}
}

func TestQuestionPairsIsolation(t *testing.T) {
	q, err := New(mustBase(t, "q"), makePairs(t, MinPairs))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := q.Pairs()
	got[0] = pair.Pair{}

	if q.Pairs()[0].IsZero() {
		t.Fatal("Pairs() must return a copy, but internal state was changed")
	}
}

func TestQuestionChangePairs(t *testing.T) {
	oldPairs := makePairs(t, MinPairs)
	newPairs := makePairs(t, MinPairs+1)
	q, err := New(mustBase(t, "q"), oldPairs)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	tests := []struct {
		name       string
		pairs      []pair.Pair
		wantErr    bool
		wantLen    int
		keepOld    bool
		wantErrMsg string
	}{
		{
			name:    "успех",
			pairs:   newPairs,
			wantErr: false,
			wantLen: len(newPairs),
		},
		{
			name:       "ошибка при пустой паре",
			pairs:      []pair.Pair{mustPair(t, "p1", "m1"), {}},
			wantErr:    true,
			wantLen:    len(newPairs),
			keepOld:    true,
			wantErrMsg: "не заполнена",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := q.ChangePairs(tt.pairs)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ChangePairs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ChangePairs() error = %v, want ErrInvalid", err)
				}
				if tt.wantErrMsg != "" && !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Fatalf("ChangePairs() error text = %q, want contain %q", err.Error(), tt.wantErrMsg)
				}
			}

			if len(q.Pairs()) != tt.wantLen {
				t.Fatalf("ChangePairs() len = %d, want %d", len(q.Pairs()), tt.wantLen)
			}
		})
	}
}

func TestQuestionClone(t *testing.T) {
	q, err := New(mustBase(t, "q"), makePairs(t, MinPairs))
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	cloned, ok := q.Clone().(*Question)
	if !ok {
		t.Fatalf("Clone() type = %T, want *Question", q.Clone())
	}
	if cloned == q {
		t.Fatal("Clone() must return new pointer")
	}
	if cloned.Base == q.Base {
		t.Fatal("Clone() must clone base pointer")
	}
	if len(cloned.Pairs()) != len(q.Pairs()) {
		t.Fatalf("Clone() pairs len = %d, want %d", len(cloned.Pairs()), len(q.Pairs()))
	}
}

func TestQuestionZeroValue(t *testing.T) {
	var q Question

	if q.Type() != question.TypeMatching {
		t.Fatalf("Type() = %v, want %v", q.Type(), question.TypeMatching)
	}
	if q.Instruction() != question.TypeMatching.DefaultInstruction() {
		t.Fatalf("Instruction() = %q, want %q", q.Instruction(), question.TypeMatching.DefaultInstruction())
	}
	if q.Pairs() != nil {
		t.Fatalf("Pairs() = %v, want nil", q.Pairs())
	}

	err := q.ChangePairs(makePairs(t, MinPairs))
	if err != nil {
		t.Fatalf("ChangePairs() zero value error = %v", err)
	}
}

func TestValidateFunctions(t *testing.T) {
	validPairs := makePairs(t, MinPairs)

	t.Run("validateRequiredBase", func(t *testing.T) {
		tests := []struct {
			name    string
			base    *base.Base
			wantErr bool
		}{
			{name: "base есть", base: mustBase(t, "q"), wantErr: false},
			{name: "base nil", base: nil, wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validateRequiredBase(tt.base)
				if (err != nil) != tt.wantErr {
					t.Fatalf("validateRequiredBase() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr && !errors.Is(err, ErrInvalid) {
					t.Fatalf("validateRequiredBase() error = %v, want ErrInvalid", err)
				}
			})
		}
	})

	t.Run("validatePairsCount", func(t *testing.T) {
		tests := []struct {
			name    string
			pairs   []pair.Pair
			wantErr bool
		}{
			{name: "min граница", pairs: makePairs(t, MinPairs), wantErr: false},
			{name: "max граница", pairs: makePairs(t, MaxPairs), wantErr: false},
			{name: "ниже min", pairs: makePairs(t, MinPairs-1), wantErr: true},
			{name: "выше max", pairs: makePairs(t, MaxPairs+1), wantErr: true},
			{name: "nil slice", pairs: nil, wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validatePairsCount(tt.pairs)
				if (err != nil) != tt.wantErr {
					t.Fatalf("validatePairsCount() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr && !errors.Is(err, ErrInvalid) {
					t.Fatalf("validatePairsCount() error = %v, want ErrInvalid", err)
				}
			})
		}
	})

	t.Run("validatePairsContainsEmpty", func(t *testing.T) {
		tests := []struct {
			name    string
			pairs   []pair.Pair
			wantErr bool
		}{
			{name: "все пары валидные", pairs: validPairs, wantErr: false},
			{name: "есть zero pair", pairs: []pair.Pair{mustPair(t, "a", "1"), {}}, wantErr: true},
			{name: "пустой slice", pairs: []pair.Pair{}, wantErr: false},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validatePairsContainsEmpty(tt.pairs)
				if (err != nil) != tt.wantErr {
					t.Fatalf("validatePairsContainsEmpty() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr && !errors.Is(err, ErrInvalid) {
					t.Fatalf("validatePairsContainsEmpty() error = %v, want ErrInvalid", err)
				}
			})
		}
	})

	t.Run("validatePairs", func(t *testing.T) {
		tests := []struct {
			name    string
			pairs   []pair.Pair
			wantErr bool
		}{
			{name: "валидный набор", pairs: validPairs, wantErr: false},
			{name: "невалидный набор по count", pairs: makePairs(t, MinPairs-1), wantErr: true},
			{name: "невалидный набор по zero pair", pairs: []pair.Pair{mustPair(t, "a", "1"), {}}, wantErr: true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := validatePairs(tt.pairs)
				if (err != nil) != tt.wantErr {
					t.Fatalf("validatePairs() error = %v, wantErr %v", err, tt.wantErr)
				}
				if tt.wantErr && !errors.Is(err, ErrInvalid) {
					t.Fatalf("validatePairs() error = %v, want ErrInvalid", err)
				}
			})
		}
	})
}
