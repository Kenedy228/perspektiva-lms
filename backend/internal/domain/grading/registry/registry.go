package registry

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/question"
)

var (
	ErrDuplicateType = fmt.Errorf("дублирующийся checker для типа вопроса")
	ErrNotFound      = fmt.Errorf("checker для типа вопроса не найден")
)

type Registry struct {
	checkers map[question.Type]grading.Checker
}

func New(checkers map[question.Type]grading.Checker) (*Registry, error) {
	for t, c := range checkers {
		if c == nil {
			return nil, fmt.Errorf("%w: тип %s — checker не может быть nil", ErrDuplicateType, t)
		}
		if !t.IsValid() {
			return nil, fmt.Errorf("%w: неизвестный тип вопроса %s", ErrDuplicateType, t)
		}
	}

	return &Registry{checkers: checkers}, nil
}

func (r *Registry) Get(t question.Type) (grading.Checker, error) {
	c, ok := r.checkers[t]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrNotFound, t)
	}

	return c, nil
}
