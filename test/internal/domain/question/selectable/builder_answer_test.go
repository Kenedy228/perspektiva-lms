package selectable_test

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/selectable"
	"github.com/google/uuid"
)

func newAnswerBuilder() *answerBuilder {
	return &answerBuilder{
		optionIDs: []uuid.UUID{},
	}
}

type answerBuilder struct {
	optionIDs []uuid.UUID
}

func (b *answerBuilder) withID(id uuid.UUID) *answerBuilder {
	b.optionIDs = append(b.optionIDs, id)
	return b
}

func (b *answerBuilder) withRandomID() *answerBuilder {
	b.optionIDs = append(b.optionIDs, uuid.New())
	return b
}

func (b *answerBuilder) build() question.Answer {
	params := selectable.AnswerParams{
		OptionIDs: b.optionIDs,
	}

	return selectable.NewAnswer(params)
}
