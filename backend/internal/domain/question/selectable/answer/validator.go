package answer

import (
	"fmt"

	"gitflic.ru/lms/backend/internal/domain/question/selectable"
	"github.com/google/uuid"
)

func validateOptionIDs(optionIDs []uuid.UUID) error {
	if err := validateOptionIDsCount(optionIDs); err != nil {
		return err
	}

	if err := validateOptionIDsContainsEmpty(optionIDs); err != nil {
		return err
	}

	if err := validateOptionIDsDuplicates(optionIDs); err != nil {
		return err
	}

	return nil
}

func validateOptionIDsCount(optionIDs []uuid.UUID) error {
	if len(optionIDs) > selectable.MaxOptionsCount {
		return fmt.Errorf(
			"%w: ответ на вопрос с выбором не может содержать больше %d выбранных опций (текущее количество - %d)",
			ErrInvalid,
			selectable.MaxOptionsCount,
			len(optionIDs),
		)
	}

	return nil
}

func validateOptionIDsContainsEmpty(optionIDs []uuid.UUID) error {
	for i := range optionIDs {
		if optionIDs[i] == uuid.Nil {
			return fmt.Errorf("%w: идентификатор выбранной опции в позиции %d не может быть пустым", ErrInvalid, i)
		}
	}

	return nil
}

func validateOptionIDsDuplicates(optionIDs []uuid.UUID) error {
	for i := range optionIDs {
		for j := i + 1; j < len(optionIDs); j++ {
			if optionIDs[i] == optionIDs[j] {
				return fmt.Errorf("%w: выбранная опция %s указана в ответе больше одного раза", ErrInvalid, optionIDs[i])
			}
		}
	}

	return nil
}
