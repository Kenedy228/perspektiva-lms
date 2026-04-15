package quiz

import (
	"errors"
	"slices"

	"github.com/google/uuid"
)

type CriteriaType string

const (
	CriteriaTypeRandom CriteriaType = "random"
	CriteriaTypeManual CriteriaType = "manual"
)

var (
	ErrEmptyQuestions = errors.New("empty questions")
	ErrNilQuestion    = errors.New("nil question found")
	ErrInvalidCount   = errors.New("invalid criteria count")
)

type Criteria interface {
	Type() CriteriaType
	QuestionCount() int
}

type RandomCriteria struct {
	count int
}

func NewRandomCriteria(count int) (RandomCriteria, error) {
	if count <= 0 {
		return RandomCriteria{}, ErrInvalidCount
	}

	return RandomCriteria{
		count: count,
	}, nil
}

func (c RandomCriteria) Count() int {
	return c.count
}

func (c RandomCriteria) QuestionCount() int {
	return c.Count()
}

func (c RandomCriteria) Type() CriteriaType {
	return CriteriaTypeRandom
}

type ManualCriteria struct {
	questionIDs []uuid.UUID
}

func NewManualCriteria(questionIDs []uuid.UUID) (ManualCriteria, error) {
	if len(questionIDs) == 0 {
		return ManualCriteria{}, ErrEmptyQuestions
	}

	if slices.Contains(questionIDs, uuid.Nil) {
		return ManualCriteria{}, ErrNilQuestion
	}

	qCopy := slices.Clone(questionIDs)

	return ManualCriteria{
		questionIDs: qCopy,
	}, nil
}

func (c ManualCriteria) QuestionIDs() []uuid.UUID {
	return slices.Clone(c.questionIDs)
}

func (c ManualCriteria) Type() CriteriaType {
	return CriteriaTypeManual
}

func (c ManualCriteria) QuestionCount() int {
	return len(c.questionIDs)
}
