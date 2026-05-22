package selectable

import (
	"gitflic.ru/lms/backend/internal/domain/question/selectable/option"
)

func countCorrectOptions(options []option.Option) int {
	count := 0

	for i := range options {
		if options[i].IsCorrect() {
			count++
		}
	}

	return count
}
