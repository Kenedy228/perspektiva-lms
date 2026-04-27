package criteria

import (
	"slices"

	"github.com/google/uuid"
)

type Manual struct {
	questionIDs []uuid.UUID
}

func NewManual(questionIDs []uuid.UUID) (Criteria, error) {
	if err := validateQuestionIDs(questionIDs); err != nil {
		return nil, err
	}

	qCopy := slices.Clone(questionIDs)

	return Manual{
		questionIDs: qCopy,
	}, nil
}

func (c Manual) QuestionIDs() []uuid.UUID {
	return slices.Clone(c.questionIDs)
}

func (c Manual) QuestionCount() int {
	return len(c.questionIDs)
}

func (c Manual) Type() Type {
	return TypeManual
}
