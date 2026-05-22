//go:build legacy
// +build legacy

package source_test

import (
	"gitflic.ru/lms/backend/internal/domain/quiz/source/criteria"
)

type criteriaFixture struct{}

func (f criteriaFixture) Type() criteria.Type {
	panic("")
}

func (f criteriaFixture) QuestionCount() int {
	panic("")
}
