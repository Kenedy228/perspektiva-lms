package attempt

import (
	"slices"

	"gitflic.ru/lms/internal/domain/question"
	"gitflic.ru/lms/internal/domain/utils"
	"github.com/google/uuid"
)

type Item struct {
	id            uuid.UUID
	questionType  question.Type
	snapshot      []byte
	studentAnswer []byte
	isCorrect     *bool
}

func NewItem(questionType question.Type, snapshot []byte) (Item, error) {
	if err := validateQuestionType(questionType); err != nil {
		return Item{}, err
	}

	if err := validateSnapshot(snapshot); err != nil {
		return Item{}, err
	}

	cSnapshot := slices.Clone(snapshot)

	id, err := utils.GenerateID()
	if err != nil {
		return Item{}, err
	}

	return Item{
		id:           id,
		questionType: questionType,
		snapshot:     cSnapshot,
	}, nil
}

func (i Item) ID() uuid.UUID {
	return i.id
}

func (i Item) QuestionType() question.Type {
	return i.questionType
}

func (i Item) Snapshot() []byte {
	return slices.Clone(i.snapshot)
}

func (i Item) StudentAnswer() []byte {
	return slices.Clone(i.studentAnswer)
}

func (i Item) IsCorrect() *bool {
	if i.isCorrect == nil {
		return nil
	}

	v := *i.isCorrect
	return &v
}

func (i Item) SetAnswer(answer []byte) error {
	if err := validateAnswer(answer); err != nil {
		return err
	}

	i.studentAnswer = answer
	return nil
}

func (i Item) IsAnswered() bool {
	return len(i.studentAnswer) > 0
}
