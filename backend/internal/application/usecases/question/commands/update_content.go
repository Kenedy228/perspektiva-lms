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
	"gitflic.ru/lms/backend/internal/domain/role"
)

// ChangeSelectableOptionsUseCase изменяет варианты вопроса с выбором.
type ChangeSelectableOptionsUseCase struct{ r questports.Repository }

// NewChangeSelectableOptionsUseCase создает ChangeSelectableOptionsUseCase.
func NewChangeSelectableOptionsUseCase(r questports.Repository) *ChangeSelectableOptionsUseCase {
	if r == nil {
		panic("question change selectable options usecase requires repository")
	}
	return &ChangeSelectableOptionsUseCase{r: r}
}

// ChangeSelectableOptionsInput описывает входные данные для изменения вариантов выбора.
type ChangeSelectableOptionsInput struct {
	ActorRole  role.Role
	QuestionID string
	Options    []SelectableOptionInput
}

// Execute изменяет варианты вопроса с выбором и сохраняет вопрос.
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
		return nil, fmt.Errorf("%w: вопрос не относится к типу с выбором", common.ErrInvalidInput)
	}
	options, err := buildSelectableOptions(in.Options)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeOptions(options); err != nil {
		return nil, fmt.Errorf("изменение вариантов вопроса с выбором: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

// ChangeSequenceOptionsUseCase изменяет варианты вопроса на последовательность.
type ChangeSequenceOptionsUseCase struct{ r questports.Repository }

// NewChangeSequenceOptionsUseCase создает ChangeSequenceOptionsUseCase.
func NewChangeSequenceOptionsUseCase(r questports.Repository) *ChangeSequenceOptionsUseCase {
	if r == nil {
		panic("question change sequence options usecase requires repository")
	}
	return &ChangeSequenceOptionsUseCase{r: r}
}

// ChangeSequenceOptionsInput описывает входные данные для изменения последовательности.
type ChangeSequenceOptionsInput struct {
	ActorRole  role.Role
	QuestionID string
	Options    []SequenceOptionInput
}

// Execute изменяет варианты вопроса на последовательность и сохраняет вопрос.
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
		return nil, fmt.Errorf("%w: вопрос не относится к типу последовательности", common.ErrInvalidInput)
	}
	options, err := buildSequenceOptions(in.Options)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeOptions(options); err != nil {
		return nil, fmt.Errorf("изменение вариантов последовательности: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

// ChangeMatchingPairsUseCase изменяет пары вопроса на сопоставление.
type ChangeMatchingPairsUseCase struct{ r questports.Repository }

// NewChangeMatchingPairsUseCase создает ChangeMatchingPairsUseCase.
func NewChangeMatchingPairsUseCase(r questports.Repository) *ChangeMatchingPairsUseCase {
	if r == nil {
		panic("question change matching pairs usecase requires repository")
	}
	return &ChangeMatchingPairsUseCase{r: r}
}

// ChangeMatchingPairsInput описывает входные данные для изменения пар сопоставления.
type ChangeMatchingPairsInput struct {
	ActorRole  role.Role
	QuestionID string
	Pairs      []MatchingPairInput
}

// Execute изменяет пары вопроса на сопоставление и сохраняет вопрос.
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
		return nil, fmt.Errorf("%w: вопрос не относится к типу сопоставления", common.ErrInvalidInput)
	}
	pairs, err := buildMatchingPairs(in.Pairs)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangePairs(pairs); err != nil {
		return nil, fmt.Errorf("изменение пар сопоставления: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}

// ChangeShortVariantsUseCase изменяет варианты короткого ответа.
type ChangeShortVariantsUseCase struct{ r questports.Repository }

// NewChangeShortVariantsUseCase создает ChangeShortVariantsUseCase.
func NewChangeShortVariantsUseCase(r questports.Repository) *ChangeShortVariantsUseCase {
	if r == nil {
		panic("question change short variants usecase requires repository")
	}
	return &ChangeShortVariantsUseCase{r: r}
}

// ChangeShortVariantsInput описывает входные данные для изменения коротких ответов.
type ChangeShortVariantsInput struct {
	ActorRole  role.Role
	QuestionID string
	Variants   []ShortVariantInput
}

// Execute изменяет варианты короткого ответа и сохраняет вопрос.
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
		return nil, fmt.Errorf("%w: вопрос не относится к типу короткого ответа", common.ErrInvalidInput)
	}
	variants, err := buildShortVariants(in.Variants)
	if err != nil {
		return nil, err
	}
	if err := cast.ChangeVariants(variants); err != nil {
		return nil, fmt.Errorf("изменение вариантов короткого ответа: %w", err)
	}
	if err := uc.r.Save(ctx, cast); err != nil {
		return nil, fmt.Errorf("сохранение вопроса: %w", err)
	}
	return &Output{ID: cast.ID().String()}, nil
}
