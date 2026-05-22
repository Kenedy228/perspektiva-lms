package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	"gitflic.ru/lms/backend/internal/domain/question"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
)

type CreateUseCase struct {
	r questports.Repository
}

func NewCreateUseCase(r questports.Repository) *CreateUseCase {
	if r == nil {
		panic("question create usecase requires repository")
	}
	return &CreateUseCase{r: r}
}

type CreateInput struct {
	ActorRole         role.Role
	Type              string
	Title             string
	Attachment        *AttachmentInput
	SelectableOptions []SelectableOptionInput
	SequenceOptions   []SequenceOptionInput
	MatchingPairs     []MatchingPairInput
	TypedBlanks       []TypedBlankInput
	ShortVariants     []ShortVariantInput
}

type Output struct {
	ID string
}

func (uc *CreateUseCase) Execute(ctx context.Context, in CreateInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}

	qType, err := question.ParseType(in.Type)
	if err != nil {
		return nil, fmt.Errorf("parse question type: %w", err)
	}

	t, err := title.New(in.Title)
	if err != nil {
		return nil, fmt.Errorf("create question title: %w", err)
	}

	q, err := createQuestion(qType, t, in)
	if err != nil {
		return nil, fmt.Errorf("create question aggregate: %w", err)
	}

	if in.Attachment != nil {
		att, err := buildAttachment(*in.Attachment)
		if err != nil {
			return nil, err
		}
		q.ChangeAttachment(att)
	}

	if err := uc.r.Save(ctx, q); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}

	return &Output{ID: q.ID().String()}, nil
}
