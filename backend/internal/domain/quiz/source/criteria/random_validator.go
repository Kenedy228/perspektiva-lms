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
		return fmt.Errorf("%w: количество вопросов должно быть положительным", ErrInvalid)
	}

	if count > maxQuestionsCount {
		return fmt.Errorf("%w: количество вопросов не должно превышать %d", ErrInvalid, maxQuestionsCount)
	}

	return nil
}
