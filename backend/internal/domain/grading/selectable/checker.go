package selectable

import (
	"gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	"gitflic.ru/lms/backend/internal/domain/question/selectable/answer"
)

// Checker проверяет ответы для вопросов с выбором вариантов.
type Checker struct{}

// New создает checker для вопросов типа selectable.
func New() Checker {
	return Checker{}
}

// Check возвращает 1, если выбран ровно полный набор правильных вариантов, иначе 0.
func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	qCast, ok := q.(*selectable.Question)
	if !ok {
		return score.Score{}, grading.ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, grading.ErrInvalidAnswerType
	}

	opts := qCast.Options()
	studentAnswers := aCast.OptionIDSet()

	if len(studentAnswers) > qCast.CorrectOptionsCount() {
		return score.New(0)
	}

	correct := 0
	for i := range opts {
		_, selected := studentAnswers[opts[i].ID()]
		if !selected {
			continue
		}

		if opts[i].IsCorrect() {
			correct++
			continue
		}

		return score.New(0)
	}

	if correct != qCast.CorrectOptionsCount() {
		return score.New(0)
	}

	return score.New(1)
}
