package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	questiontitle "gitflic.ru/lms/backend/internal/domain/question/base/title"
	"gitflic.ru/lms/backend/internal/domain/role"
)

// ChangeTitleUseCase изменяет заголовок существующего вопроса.
type ChangeTitleUseCase struct {
	r questports.Repository
}

// NewChangeTitleUseCase создает ChangeTitleUseCase.
func NewChangeTitleUseCase(r questports.Repository) *ChangeTitleUseCase {
	if r == nil {
		panic("question change title usecase requires repository")
	}
	return &ChangeTitleUseCase{r: r}
}

// ChangeTitleInput описывает входные данные для изменения заголовка вопроса.
type ChangeTitleInput struct {
	ActorRole  role.Role
	QuestionID string
	Title      string
}

// Execute изменяет заголовок вопроса и сохраняет его в репозитории.
func (uc *ChangeTitleUseCase) Execute(ctx context.Context, in ChangeTitleInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	t, err := questiontitle.New(in.Title)
	if err != nil {
		return nil, fmt.Errorf("создание заголовка вопроса: %w", err)
	}

	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}

	if err := q.ChangeTitle(t); err != nil {
		return nil, fmt.Errorf("изменение заголовка вопроса: %w", err)
	}
	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}
