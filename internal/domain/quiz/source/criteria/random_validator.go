package criteria

import "fmt"

func validateQuestionCount(count int) error {
	if err := validateQuestionCountBoundaries(count); err != nil {
		return err
	}

	return nil
}

func validateQuestionCountBoundaries(count int) error {
	if count <= 0 {
		return fmt.Errorf("%w, детали: размер выборки должен быть положительным числом", ErrInvalid)
	}

	if count > maxQuestionsCount {
		return fmt.Errorf("%w, детали: размер выборки не должен превышать %d вопросов", ErrInvalid, maxQuestionsCount)
	}

	return nil
}
