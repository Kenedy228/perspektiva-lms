package criteria

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestNewManual(t *testing.T) {
	ids := []uuid.UUID{uuid.New(), uuid.New()}

	tests := []struct {
		name    string
		ids     []uuid.UUID
		wantErr bool
	}{
		{name: "один вопрос", ids: []uuid.UUID{uuid.New()}},
		{name: "несколько вопросов", ids: ids},
		{name: "пустой список", ids: nil, wantErr: true},
		{name: "пустой срез", ids: []uuid.UUID{}, wantErr: true},
		{name: "содержит nil UUID", ids: []uuid.UUID{uuid.Nil}, wantErr: true},
		{name: "содержит дубликаты", ids: []uuid.UUID{ids[0], ids[0]}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewManual(tt.ids)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got.Type() != TypeManual {
				t.Fatalf("ожидался тип TypeManual, получено %v", got.Type())
			}
			if got.QuestionCount() != len(tt.ids) {
				t.Fatalf("ожидалось QuestionCount()=%d, получено %d", len(tt.ids), got.QuestionCount())
			}
		})
	}
}

func TestNewRandom(t *testing.T) {
	tests := []struct {
		name    string
		count   int
		wantErr bool
	}{
		{name: "один вопрос", count: 1},
		{name: "несколько вопросов", count: 10},
		{name: "максимальное значение", count: maxQuestionsCount},
		{name: "ноль", count: 0, wantErr: true},
		{name: "отрицательное значение", count: -1, wantErr: true},
		{name: "превышение максимума", count: maxQuestionsCount + 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRandom(tt.count)
			if tt.wantErr {
				if !errors.Is(err, ErrInvalid) {
					t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("неожиданная ошибка: %v", err)
			}
			if got.Type() != TypeRandom {
				t.Fatalf("ожидался тип TypeRandom, получено %v", got.Type())
			}
			if got.QuestionCount() != tt.count {
				t.Fatalf("ожидалось QuestionCount()=%d, получено %d", tt.count, got.QuestionCount())
			}
		})
	}
}

func TestManual_QuestionsAreIsolated(t *testing.T) {
	original := []uuid.UUID{uuid.New()}
	c, err := NewManual(original)
	if err != nil {
		t.Fatalf("create manual criteria: %v", err)
	}

	m := c.(Manual)
	got := m.QuestionIDs()
	got[0] = uuid.Nil

	if m.QuestionIDs()[0] == uuid.Nil {
		t.Fatal("внешнее изменение не должно влиять на внутренние данные критерия")
	}
}
