package bank

import (
	"fmt"
	"strings"

	"gitflic.ru/lms/internal/domain/shared/duplicate"
	"github.com/google/uuid"
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("%w, детали: заголовок банка должен содержать хотя бы один непробельный символ", ErrInvalid)
	}

	return nil
}

func validateQuestionsForAdding(questions []uuid.UUID, add []uuid.UUID) error {
	totalLen := len(questions) + len(add)

	if totalLen > maxQuestions {
		return fmt.Errorf("%w, детали: количество вопросов превысит лимит %d вопросов", ErrInvalid, maxQuestions)
	}

	for i := range add {
		if add[i] == uuid.Nil {
			return fmt.Errorf("%w, детали: добавляемые вопросы содержат несуществующий вопрос", ErrInvalid)
		}
	}

	cQuestions := make([]uuid.UUID, 0, totalLen)
	cQuestions = append(cQuestions, questions...)
	cQuestions = append(cQuestions, add...)

	if has := duplicate.FindUUID(cQuestions); has {
		return fmt.Errorf("%w, детали: один из добавляемых вопросов либо уже существует в банке вопросов, либо дублирует вопрос в добавляемом списке", ErrInvalid)
	}

	return nil
}
