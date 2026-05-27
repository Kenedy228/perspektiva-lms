package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	questiontitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/role"
)

// CreateUseCase создает новый вопрос выбранного типа.
type CreateUseCase struct {
	r questports.Repository
}

// NewCreateUseCase создает CreateUseCase.
func NewCreateUseCase(r questports.Repository) *CreateUseCase {
	if r == nil {
		panic("question create usecase requires repository")
	}
	return &CreateUseCase{r: r}
}

// CreateInput описывает входные данные для создания вопроса.
type CreateInput struct {
	ActorRole         role.Role
	Type              string
	Title             string
	SelectableOptions []SelectableOptionInput
	SequenceOptions   []SequenceOptionInput
	MatchingPairs     []MatchingPairInput
	ShortVariants     []ShortVariantInput
}

// Output содержит результат выполнения команды.
type Output struct {
	ID string
}

// Execute создает вопрос и сохраняет его в репозитории.
func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	qType := question.Type(in.Type)
	if !qType.IsValid() {
		return nil, fmt.Errorf("%w: некорректный тип вопроса", common.ErrInvalidInput)
	}

	t, err := questiontitle.New(in.Title)
	if err != nil {
		return nil, fmt.Errorf("создание заголовка вопроса: %w", err)
	}

	q, err := createQuestion(qType, t, in)
	if err != nil {
		return nil, fmt.Errorf("создание агрегата вопроса: %w", err)
	}

	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}
