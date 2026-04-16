package manual

import (
	"slices"

	"gitflic.ru/lms/internal/domain/quiz/criteria"
	"github.com/google/uuid"
)

type ManualCriteria struct {
	questionIDs []uuid.UUID
}

func NewManualCriteria(params Params) (ManualCriteria, error) {
	if err := validateQuestionIDs(params.QuestionIDs); err != nil {
		return ManualCriteria{}, err
	}

	qCopy := slices.Clone(params.QuestionIDs)

	return ManualCriteria{
		questionIDs: qCopy,
	}, nil
}

func (c ManualCriteria) QuestionIDs() []uuid.UUID {
	return slices.Clone(c.questionIDs)
}

func (c ManualCriteria) Type() criteria.CriteriaType {
	return criteria.CriteriaTypeManual
}

func (c ManualCriteria) QuestionCount() int {
	return len(c.questionIDs)
}
