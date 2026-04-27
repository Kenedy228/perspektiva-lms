package matching_test

import (
	"testing"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/question/matching"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func newQuestionBuilder() *questionBuilder {
	return &questionBuilder{
		text:    "text",
		imageID: uuid.Nil,
		pairs:   []matching.Pair{},
	}
}

type questionBuilder struct {
	text    string
	imageID uuid.UUID
	pairs   []matching.Pair
}

func (b *questionBuilder) withText(text string) *questionBuilder {
	b.text = text
	return b
}

func (b *questionBuilder) withImageID(id uuid.UUID) *questionBuilder {
	b.imageID = id
	return b
}

func (b *questionBuilder) withPairParam(prompt, val string, cType question.ContentType) *questionBuilder {
	pair := newPairBuilder().withPrompt(prompt).withContent(cType, val).buildNoTest()

	b.pairs = append(b.pairs, pair)
	return b
}

func (b *questionBuilder) withPairParamAsText(prompt, val string) *questionBuilder {
	return b.withPairParam(prompt, val, question.ContentTypeText)
}

func (b *questionBuilder) build(t *testing.T, wantErr error) question.Question {
	t.Helper()

	params := matching.Params{
		Text:    b.text,
		ImageID: b.imageID,
		Pairs:   b.pairs,
	}

	q, err := matching.New(params)

	assert.ErrorIs(t, err, wantErr)
	return q
}
