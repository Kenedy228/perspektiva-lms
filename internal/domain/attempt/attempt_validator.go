package attempt

import (
	"errors"
	"strings"

	"gitflic.ru/lms/internal/domain/question"
	"github.com/google/uuid"
)

var (
	ErrInvalidQuestionType = errors.New("invalid question type")
	ErrEmptySnapshot       = errors.New("empty snapshot")
	ErrNilEnrollmentID     = errors.New("nil enrollment id")
	ErrNilQuizID           = errors.New("nil quiz id")
	ErrNilCourseID         = errors.New("nil course id")
	ErrEmptyCourseTitle    = errors.New("empty course title")
	ErrEmptyQuizTitle      = errors.New("empty quiz title")
)

func validateQuestionType(questionType question.Type) error {
	if !questionType.IsValid() {
		return ErrInvalidQuestionType
	}

	return nil
}

func validateSnapshot(snapshot []byte) error {
	if len(snapshot) == 0 {
		return ErrEmptySnapshot
	}

	return nil
}

func validateEnrollmentID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilEnrollmentID
	}

	return nil
}

func validateQuizID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilQuizID
	}

	return nil
}

func validateCourseID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrNilCourseID
	}

	return nil
}

func validateQuizTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyQuizTitle
	}

	return nil
}

func validateCourseTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyCourseTitle
	}

	return nil
}
