package grading

import (
	"context"
	"errors"
	"fmt"

	questionports "gitflic.ru/lms/backend/internal/application/ports/question"
	domaingrading "gitflic.ru/lms/backend/internal/domain/grading"
	"gitflic.ru/lms/backend/internal/domain/grading/score"
	"gitflic.ru/lms/backend/internal/domain/question"
	"github.com/google/uuid"
)

var (
	// ErrInvalidInput возвращается при некорректных входных данных usecase-а.
	ErrInvalidInput = errors.New("некорректные данные для оценки ответа")
	// ErrUnsupportedChecker возвращается, когда для типа вопроса не найден подходящий checker.
	ErrUnsupportedChecker = errors.New("проверяющий для типа вопроса не найден")
)

// GradeUseCase оркестрирует сценарий оценки ответа пользователя.
type GradeUseCase struct {
	r        questionports.Repository
	checkers []domaingrading.Checker
}

// NewGradeUseCase создает usecase оценки ответа и требует repository и минимум один checker.
func NewGradeUseCase(r questionports.Repository, checkers ...domaingrading.Checker) *GradeUseCase {
	if r == nil {
		panic("question grade usecase requires repository")
	}
	if len(checkers) == 0 {
		panic("question grade usecase requires checkers")
	}

	return &GradeUseCase{
		r:        r,
		checkers: append([]domaingrading.Checker(nil), checkers...),
	}
}

// GradeInput описывает входные данные для оценки ответа.
type GradeInput struct {
	QuestionID string
	Answer     question.Answer
}

// GradeOutput содержит вычисленный score ответа.
type GradeOutput struct {
	Score score.Score
}

// Execute загружает вопрос, выбирает checker по типу и возвращает оценку ответа.
func (uc *GradeUseCase) Execute(ctx context.Context, in GradeInput) (*GradeOutput, error) {
	if err := validateGradeInput(in); err != nil {
		return nil, err
	}

	q, checker, err := uc.loadQuestionAndChecker(ctx, in.QuestionID)
	if err != nil {
		return nil, err
	}

	s, err := checker.Check(q, in.Answer)
	if err != nil {
		return nil, fmt.Errorf("проверка ответа: %w", err)
	}

	return &GradeOutput{Score: s}, nil
}

func (uc *GradeUseCase) loadQuestionAndChecker(ctx context.Context, questionID string) (question.Question, domaingrading.Checker, error) {
	id, err := uuid.Parse(questionID)
	if err != nil {
		return nil, nil, fmt.Errorf("разбор идентификатора вопроса: %w", err)
	}

	q, err := uc.r.FindByID(ctx, id)
	if err != nil {
		return nil, nil, fmt.Errorf("поиск вопроса: %w", err)
	}

	for i := range uc.checkers {
		if uc.checkers[i].Supports(q.Type()) {
			return q, uc.checkers[i], nil
		}
	}

	return nil, nil, fmt.Errorf("%w: тип вопроса %s", ErrUnsupportedChecker, q.Type())
}

func validateGradeInput(in GradeInput) error {
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

// ValidateAnswerUseCase проверяет корректность ответа через сценарий оценки без возврата score.
type ValidateAnswerUseCase struct {
	grade *GradeUseCase
}

// NewValidateAnswerUseCase создает usecase валидации ответа на базе GradeUseCase.
func NewValidateAnswerUseCase(r questionports.Repository, checkers ...domaingrading.Checker) *ValidateAnswerUseCase {
	return &ValidateAnswerUseCase{grade: NewGradeUseCase(r, checkers...)}
}

// ValidateAnswerInput описывает входные данные для валидации ответа.
type ValidateAnswerInput struct {
	QuestionID string
	Answer     question.Answer
}

// Execute делегирует проверку в GradeUseCase.Execute без повторной загрузки вопроса.
func (uc *ValidateAnswerUseCase) Execute(ctx context.Context, in ValidateAnswerInput) error {
	_, err := uc.grade.Execute(ctx, GradeInput{QuestionID: in.QuestionID, Answer: in.Answer})
	return err
}
