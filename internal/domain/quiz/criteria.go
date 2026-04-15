package quiz

import "errors"

type CriteriaType string

const (
	CriteriaTypeRandom CriteriaType = "random"
	CriteriaTypeManual CriteriaType = "manual"
)

type Criteria interface {
	Type() CriteriaType
}

type RandomCriteria struct {
	count int
}

var (
	ErrInvalidCount = errors.New("invalid criteria count")
)

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

func (c RandomCriteria) Type() CriteriaType {
	return CriteriaTypeRandom
}
