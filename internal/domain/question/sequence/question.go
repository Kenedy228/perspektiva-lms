package sequence

import (
	"math/rand/v2"
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/base"
)

type Question struct {
	*base.Base
	elements []Element
}

func New(params Params) (question.Question, error) {
	base, err := base.New(params.baseParams())
	if err != nil {
		return nil, err
	}

	if err := validateElements(params.Elements); err != nil {
		return nil, err
	}

	cElements := slices.Clone(params.Elements)

	return &Question{
		Base:     base,
		elements: cElements,
	}, nil
}

func (q *Question) Instruction() string {
	return q.Type().DefaultInstruction()
}

func (q *Question) Elements() []Element {
	return slices.Clone(q.elements)
}

func (q *Question) ShuffledElements() []Element {
	cElements := slices.Clone(q.elements)
	rand.Shuffle(len(cElements), func(i, j int) {
		cElements[i], cElements[j] = cElements[j], cElements[i]
	})

	return cElements
}

func (q *Question) Type() question.Type {
	return question.TypeSequence
}

func (q *Question) UpdateElements(elements []Element) error {
	if err := validateElements(elements); err != nil {
		return err
	}

	cElements := slices.Clone(elements)

	q.elements = cElements
	q.Touch()
	return nil
}

func (q *Question) CheckAnswer(answer question.Answer) bool {
	cast, ok := answer.(Answer)
	if !ok {
		return false
	}

	studentAnswers := cast.IDs()

	if len(studentAnswers) != len(q.elements) {
		return false
	}

	for i := range q.elements {
		if q.elements[i].ID() != studentAnswers[i] {
			return false
		}
	}

	return true
}

func (q *Question) Clone() question.Question {
	return &Question{
		Base:     q.Base.Clone(),
		elements: q.Elements(),
	}
}
