package sequence_test

import (
	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/sequence"
	"github.com/google/uuid"
)

func newAnswerBuilder() *answerBuilder {
	return &answerBuilder{
		elementIDs: []uuid.UUID{},
	}
}

type answerBuilder struct {
	elementIDs []uuid.UUID
}

func (b *answerBuilder) withID(id uuid.UUID) *answerBuilder {
	b.elementIDs = append(b.elementIDs, id)
	return b
}

func (b *answerBuilder) withRandomID() *answerBuilder {
	b.elementIDs = append(b.elementIDs, uuid.New())
	return b
}

func (b *answerBuilder) build() question.Answer {
	params := sequence.AnswerParams{
		ElementIDs: b.elementIDs,
	}

	return sequence.NewAnswer(params)
}
