package attempt

import (
	"fmt"

	"gitflic.ru/lms/internal/domain/question"
)

func validateSnapshot(snapshot question.Question) error {
	if snapshot == nil {
		return fmt.Errorf("%w, детали: переданный снимок вопроса не существует (nil)", ErrInvalidItem)
	}

	return nil
}
