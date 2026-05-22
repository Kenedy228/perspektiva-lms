package commands

import (
	"context"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	"gitflic.ru/lms/backend/internal/application/usecases/question/common"
	qmatching "gitflic.ru/lms/backend/internal/domain/question/matching"
	qselectable "gitflic.ru/lms/backend/internal/domain/question/selectable"
	qsequence "gitflic.ru/lms/backend/internal/domain/question/sequence"
	qshort "gitflic.ru/lms/backend/internal/domain/question/short"
	qtyped "gitflic.ru/lms/backend/internal/domain/question/typed"
	"gitflic.ru/lms/backend/internal/domain/role"
	"gitflic.ru/lms/backend/internal/domain/shared/title"
)

type ChangeSelectableOptionsUseCase struct{ r questports.Repository }

func NewChangeSelectableOptionsUseCase(r questports.Repository) *ChangeSelectableOptionsUseCase {
	if r == nil {
		panic("question change selectable options usecase requires repository")
	}
	return &ChangeSelectableOptionsUseCase{r: r}
}

type ChangeSelectableOptionsInput struct {
	ActorRole  role.Role
	QuestionID string
	Options    []SelectableOptionInput
}

func (uc *ChangeSelectableOptionsUseCase) Execute(ctx context.Context, in ChangeSelectableOptionsInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}
	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}
	cast, ok := q.(*qselectable.Question)
	if !ok {
		return nil, fmt.Errorf("%w: question is not selectable", common.ErrInvalidInput)
	}
	options, err := buildSelectableOptions(in.Options)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeOptions(options); err != nil {
		return nil, fmt.Errorf("change selectable options: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

type ChangeSequenceOptionsUseCase struct{ r questports.Repository }

func NewChangeSequenceOptionsUseCase(r questports.Repository) *ChangeSequenceOptionsUseCase {
	if r == nil {
		panic("question change sequence options usecase requires repository")
	}
	return &ChangeSequenceOptionsUseCase{r: r}
}

type ChangeSequenceOptionsInput struct {
	ActorRole  role.Role
	QuestionID string
	Options    []SequenceOptionInput
}

func (uc *ChangeSequenceOptionsUseCase) Execute(ctx context.Context, in ChangeSequenceOptionsInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}
	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}
	cast, ok := q.(*qsequence.Question)
	if !ok {
		return nil, fmt.Errorf("%w: question is not sequence", common.ErrInvalidInput)
	}
	options, err := buildSequenceOptions(in.Options)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeOptions(options); err != nil {
		return nil, fmt.Errorf("change sequence options: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

type ChangeMatchingPairsUseCase struct{ r questports.Repository }

func NewChangeMatchingPairsUseCase(r questports.Repository) *ChangeMatchingPairsUseCase {
	if r == nil {
		panic("question change matching pairs usecase requires repository")
	}
	return &ChangeMatchingPairsUseCase{r: r}
}

type ChangeMatchingPairsInput struct {
	ActorRole  role.Role
	QuestionID string
	Pairs      []MatchingPairInput
}

func (uc *ChangeMatchingPairsUseCase) Execute(ctx context.Context, in ChangeMatchingPairsInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}
	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}
	cast, ok := q.(*qmatching.Question)
	if !ok {
		return nil, fmt.Errorf("%w: question is not matching", common.ErrInvalidInput)
	}
	pairs, err := buildMatchingPairs(in.Pairs)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangePairs(pairs); err != nil {
		return nil, fmt.Errorf("change matching pairs: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

type ChangeTypedContentUseCase struct{ r questports.Repository }

func NewChangeTypedContentUseCase(r questports.Repository) *ChangeTypedContentUseCase {
	if r == nil {
		panic("question change typed content usecase requires repository")
	}
	return &ChangeTypedContentUseCase{r: r}
}

type ChangeTypedContentInput struct {
	ActorRole  role.Role
	QuestionID string
	Title      string
	Blanks     []TypedBlankInput
}

func (uc *ChangeTypedContentUseCase) Execute(ctx context.Context, in ChangeTypedContentInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}
	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}
	cast, ok := q.(*qtyped.Question)
	if !ok {
		return nil, fmt.Errorf("%w: question is not typed", common.ErrInvalidInput)
	}
	t, err := title.New(in.Title)
	if err != nil {
		return nil, fmt.Errorf("create typed title: %w", err)
	}
	blanks, err := buildTypedBlanks(in.Blanks)
	if err != nil {
		return nil, err
	}
	if err := cast.ReplaceContent(t, blanks); err != nil {
		return nil, fmt.Errorf("change typed content: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

type ChangeShortVariantsUseCase struct{ r questports.Repository }

func NewChangeShortVariantsUseCase(r questports.Repository) *ChangeShortVariantsUseCase {
	if r == nil {
		panic("question change short variants usecase requires repository")
	}
	return &ChangeShortVariantsUseCase{r: r}
}

type ChangeShortVariantsInput struct {
	ActorRole  role.Role
	QuestionID string
	Variants   []ShortVariantInput
}

func (uc *ChangeShortVariantsUseCase) Execute(ctx context.Context, in ChangeShortVariantsInput) (*Output, error) {
	if err := common.RequireAuthor(in.ActorRole); err != nil {
		return nil, err
	}
	q, err := loadQuestion(ctx, uc.r, in.QuestionID)
	if err != nil {
		return nil, err
	}
	cast, ok := q.(*qshort.Question)
	if !ok {
		return nil, fmt.Errorf("%w: question is not short", common.ErrInvalidInput)
	}
	variants, err := buildShortVariants(in.Variants)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeVariants(variants); err != nil {
		return nil, fmt.Errorf("change short variants: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("save question: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}
