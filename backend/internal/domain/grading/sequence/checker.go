package sequence

import (
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/question/sequence"
	"gitflic.ru/lms/backend/internal/domain/question/sequence/answer"
	"github.com/google/uuid"
)

// Checker проверяет ответы для вопросов с упорядочиванием последовательности.
type Checker struct{}

// New создает checker для вопросов типа sequence.
func New() Checker {
	return Checker{}
}

// Check возвращает 1 только при полном совпадении длины и порядка элементов.
func (c Checker) Check(q question.Question, a question.Answer) (score.Score, error) {
	qCast, ok := q.(*sequence.Question)
	if !ok {
		return score.Score{}, ErrInvalidQuestionType
	}

	aCast, ok := a.(answer.Answer)
	if !ok {
		return score.Score{}, ErrInvalidAnswerType
	}

	seq := qCast.Options()
	studentAnswers := aCast.OptionIDs()

	if len(seq) != len(studentAnswers) {
		return score.New(0)
	}

	for i := range seq {
		if optionIDFromValue(seq[i].Value()) != studentAnswers[i].ID() {
			return score.New(0)
		}
	}

	return score.New(1)
}

// Supports сообщает, поддерживается ли тип вопроса checker-ом.
func (c Checker) Supports(qType question.Type) bool {
	return qType == question.TypeSequence
}

func optionIDFromValue(value string) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceOID, []byte(value))
}
