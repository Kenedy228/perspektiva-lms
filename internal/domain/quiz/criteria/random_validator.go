package criteria

import "fmt"

func validateQuestionCount(count int) error {
	if count <= 0 {
		return fmt.Errorf("%w, детали: размер выборки должен быть положительным числом", ErrInvalidQuestionCount)
	}

	if count > maxQuestions {
		return fmt.Errorf("%w, детали: размер выборки не должен превышать %d вопросов", ErrInvalidQuestionCount, maxQuestions)
	}

	return nil
}
