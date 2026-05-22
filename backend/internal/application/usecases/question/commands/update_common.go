package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
)

type ChangeTitleUseCase struct {
	r questports.Repository
}

func NewChangeTitleUseCase(r questports.Repository) *ChangeTitleUseCase {
	if r == nil {
		panic("question change title usecase requires repository")
	}
	return &ChangeTitleUseCase{r: r}
}

type ChangeTitleInput struct {
	ActorRole  role.Role
	QuestionID string
	Title      string
}

func (uc *ChangeTitleUseCase) Execute(ctx context.Context, in ChangeTitleInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	t, err := title.New(in.Title)
	if err != nil {
		return nil, fmt.Errorf("create question title: %w", err)
	}

	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}

	q.ChangeTitle(t)
	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}

type ChangeAttachmentUseCase struct {
	r questports.Repository
}

func NewChangeAttachmentUseCase(r questports.Repository) *ChangeAttachmentUseCase {
	if r == nil {
		panic("question change attachment usecase requires repository")
	}
	return &ChangeAttachmentUseCase{r: r}
}

type ChangeAttachmentInput struct {
	ActorRole  role.Role
	QuestionID string
	Attachment AttachmentInput
}

func (uc *ChangeAttachmentUseCase) Execute(ctx context.Context, in ChangeAttachmentInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	att, err := buildAttachment(in.Attachment)
	if err != nil {
		return nil, err
	}

	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}

	q.ChangeAttachment(att)
	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}

type RemoveAttachmentUseCase struct {
	r questports.Repository
}

func NewRemoveAttachmentUseCase(r questports.Repository) *RemoveAttachmentUseCase {
	if r == nil {
		panic("question remove attachment usecase requires repository")
	}
	return &RemoveAttachmentUseCase{r: r}
}

type RemoveAttachmentInput struct {
	ActorRole  role.Role
	QuestionID string
}

func (uc *RemoveAttachmentUseCase) Execute(ctx context.Context, in RemoveAttachmentInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}

	q.RemoveAttachment()
	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}
