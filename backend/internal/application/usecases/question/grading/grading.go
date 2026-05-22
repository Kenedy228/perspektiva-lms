package grading

import (
	"context"
	"errors"
	"fmt"

	questports "gitflic.ru/lms/backend/internal/application/ports/question"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

var (
	ErrInvalidInput       = errors.New("question grading invalid input")
	ErrUnsupportedChecker = errors.New("question grading checker not found")
)

type GradeUseCase struct {
	r        questports.Repository
	checkers []domaingrading.Checker
}

func NewGradeUseCase(r questports.Repository, checkers ...domaingrading.Checker) *GradeUseCase {
	if r == nil {
		panic("question grade usecase requires repository")
	}
	if len(checkers) == 0 {
		panic("question grade usecase requires checkers")
	}
	return &GradeUseCase{r: r, checkers: append([]domaingrading.Checker(nil), checkers...)}
}

type GradeInput struct {
	QuestionID string
	Answer     question.Answer
}

type GradeOutput struct {
	Score score.Score
}

func (uc *GradeUseCase) Execute(ctx context.Context, in GradeInput) (*GradeOutput, error) {
	q, checker, err := uc.load(ctx, in.QuestionID)
	if err != nil {
		return nil, err
	}

	s, err := checker.Check(q, in.Answer)
	if err != nil {
		return nil, fmt.Errorf("check answer: %w", err)
	}

	return &GradeOutput{Score: s}, nil
}

func (uc *GradeUseCase) load(ctx context.Context, questionID string) (question.Question, domaingrading.Checker, error) {
	id, err := uuid.Parse(questionID)
	if err != nil {
		return nil, nil, fmt.Errorf("parse question id: %w", err)
	}
	if id == uuid.Nil {
		return nil, nil, fmt.Errorf("%w: question id is required", ErrInvalidInput)
	}

	q, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("find question: %w", err)
	}

	for i := range uc.checkers {
		if uc.checkers[i].Supports(q.Type()) {
			return q, uc.checkers[i], nil
		}
	}

	return nil, nil, fmt.Errorf("%w: %s", ErrUnsupportedChecker, q.Type())
}

type ValidateAnswerUseCase struct {
	grade *GradeUseCase
}

func NewValidateAnswerUseCase(r questports.Repository, checkers ...domaingrading.Checker) *ValidateAnswerUseCase {
	return &ValidateAnswerUseCase{
		grade: NewGradeUseCase(r, checkers...),
	}
}

type ValidateAnswerInput struct {
	QuestionID string
	Answer     question.Answer
}

func (uc *ValidateAnswerUseCase) Execute(ctx context.Context, in ValidateAnswerInput) error {
	_, _, err := uc.grade.load(ctx, in.QuestionID)
	if err != nil {
		return err
	}

	_, err = uc.grade.Execute(ctx, GradeInput{
		QuestionID: in.QuestionID,
		Answer:     in.Answer,
	})
	return err
}
