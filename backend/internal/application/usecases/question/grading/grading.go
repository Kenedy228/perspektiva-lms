package grading

import (
	"context"
	"errors"
	"fmt"

	questionports "gitflic.ru/lms/backend/internal/application/ports/question"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/registry"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

var ErrInvalidInput = errors.New("некорректные данные для оценки ответа")

type GradeUseCase struct {
	r          questionports.Repository
	registry   *registry.Registry
	validators map[question.Type]domaingrading.AnswerValidator
}

func NewGradeUseCase(
	r questionports.Repository,
	reg *registry.Registry,
	validators map[question.Type]domaingrading.AnswerValidator,
) *GradeUseCase {
	if r == nil {
		panic("question grade usecase requires repository")
	}
	if reg == nil {
		panic("question grade usecase requires checker registry")
	}
	if len(validators) == 0 {
		panic("question grade usecase requires validators")
	}

	return &GradeUseCase{
		r:          r,
		registry:   reg,
		validators: validators,
	}
}

type GradeInput struct {
	QuestionID string
	Answer     question.Answer
}

type GradeOutput struct {
	Score score.Score
}

func (uc *GradeUseCase) Execute(ctx context.Context, in GradeInput) (*GradeOutput, error) {
	if err := validateInput(in); err != nil {
		return nil, err
	}

	q, err := uc.loadQuestion(ctx, in.QuestionID)
	if err != nil {
		return nil, err
	}

	v, ok := uc.validators[q.Type()]
	if !ok {
		return nil, fmt.Errorf("валидатор для типа вопроса %s не найден", q.Type())
	}

	if err := v.Validate(q, in.Answer); err != nil {
		return nil, fmt.Errorf("валидация ответа: %w", err)
	}

	checker, err := uc.registry.Get(q.Type())
	if err != nil {
		return nil, fmt.Errorf("поиск проверяющего: %w", err)
	}

	s, err := checker.Check(q, in.Answer)
	if err != nil {
		return nil, fmt.Errorf("проверка ответа: %w", err)
	}

	return &GradeOutput{Score: s}, nil
}

func (uc *GradeUseCase) loadQuestion(ctx context.Context, questionID string) (question.Question, error) {
	id, err := uuid.Parse(questionID)
	if err != nil {
		return nil, fmt.Errorf("разбор идентификатора вопроса: %w", err)
	}

	q, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("поиск вопроса: %w", err)
	}

	return q, nil
}

type ValidateAnswerUseCase struct {
	r          questionports.Repository
	validators map[question.Type]domaingrading.AnswerValidator
}

func NewValidateAnswerUseCase(
	r questionports.Repository,
	validators map[question.Type]domaingrading.AnswerValidator,
) *ValidateAnswerUseCase {
	if r == nil {
		panic("validate answer usecase requires repository")
	}
	if len(validators) == 0 {
		panic("validate answer usecase requires validators")
	}

	return &ValidateAnswerUseCase{
		r:          r,
		validators: validators,
	}
}

type ValidateAnswerInput struct {
	QuestionID string
	Answer     question.Answer
}

func (uc *ValidateAnswerUseCase) Execute(ctx context.Context, in ValidateAnswerInput) error {
	if err := validateInput(GradeInput{QuestionID: in.QuestionID, Answer: in.Answer}); err != nil {
		return err
	}

	q, err := uc.loadQuestion(ctx, in.QuestionID)
	if err != nil {
		return err
	}

	v, ok := uc.validators[q.Type()]
	if !ok {
		return fmt.Errorf("валидатор для типа вопроса %s не найден", q.Type())
	}

	if err := v.Validate(q, in.Answer); err != nil {
		return fmt.Errorf("валидация ответа: %w", err)
	}

	return nil
}

func (uc *ValidateAnswerUseCase) loadQuestion(ctx context.Context, questionID string) (question.Question, error) {
	id, err := uuid.Parse(questionID)
	if err != nil {
		return nil, fmt.Errorf("разбор идентификатора вопроса: %w", err)
	}

	q, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("поиск вопроса: %w", err)
	}

	return q, nil
}

func validateInput(in GradeInput) error {
	if in.QuestionID == "" {
		return fmt.Errorf("%w: идентификатор вопроса обязателен", ErrInvalidInput)
	}

	id, err := uuid.Parse(in.QuestionID)
	if err != nil {
		return fmt.Errorf("%w: идентификатор вопроса имеет некорректный формат", ErrInvalidInput)
	}

	if id == uuid.Nil {
		return fmt.Errorf("%w: идентификатор вопроса не может быть uuid.Nil", ErrInvalidInput)
	}

	if in.Answer == nil {
		return fmt.Errorf("%w: ответ обязателен", ErrInvalidInput)
	}

	return nil
}
