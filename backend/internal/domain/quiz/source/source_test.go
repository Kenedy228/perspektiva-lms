package source

import (
	"errors"
	"testing"

	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
	"github.com/google/uuid"
)

func TestNewSource_Valid(t *testing.T) {
	bankID := uuid.New()
	c := mustRandom(t, 5)

	s, err := NewSource(bankID, c)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if s.BankID() != bankID {
		t.Fatal("BankID не совпадает")
	}
	if s.Criteria() == nil {
		t.Fatal("Criteria не должны быть nil")
	}
	if s.IsZero() {
		t.Fatal("ожидалось IsZero()=false для корректного источника")
	}
}

func TestNewSource_RejectsNilBankID(t *testing.T) {
	_, err := NewSource(uuid.Nil, mustRandom(t, 3))
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
	}
}

func TestNewSource_RejectsNilCriteria(t *testing.T) {
	_, err := NewSource(uuid.New(), nil)
	if !errors.Is(err, ErrInvalid) {
		t.Fatalf("ожидалась ошибка ErrInvalid, получено: %v", err)
	}
}

func TestSource_IsZero(t *testing.T) {
	var zero Source
	if !zero.IsZero() {
		t.Fatal("ожидалось IsZero()=true для нулевого значения")
	}
}

func mustRandom(t *testing.T, count int) criteria.Criteria {
	t.Helper()
	c, err := criteria.NewRandom(count)
	if err != nil {
		t.Fatalf("create random criteria: %v", err)
	}
	return c
}
