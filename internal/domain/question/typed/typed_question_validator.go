package typed

import (
	"errors"
	"regexp"

	"gitflic.ru/lms/internal/domain/question/option"
	"gitflic.ru/lms/internal/domain/shared/duplicate"
)

var (
	placeholderRegexp = regexp.MustCompile(`\{\{\S+?\}\}`)
)

// NOTE: placeholder errors
var (
	ErrDuplicatePlaceholder  = errors.New("duplicate placeholder")
	ErrTooManyPlaceholders   = errors.New("too many placeholders")
	ErrNotEnoughPlaceholders = errors.New("not enough placeholders")
	ErrCountMismatch         = errors.New("count mismatch")
	ErrInvalidBlanks         = errors.New("blanks doesn't contain placeholder")
)

// NOTE: blank errors
var (
	ErrInvalidPlaceholder  = errors.New("invalid placeholder")
	ErrEmptyAnswers        = errors.New("empty answers")
	ErrTooManyAnswers      = errors.New("too many answers")
	ErrInvalidAnswerFormat = errors.New("invalid answer format")
)

func validatePlaceholders(text string, blanks []BlankParams) error {
	tPlaceholders := placeholderRegexp.FindAllString(text, -1)

	if dupl := duplicate.FindAllComparable(tPlaceholders); len(dupl) != 0 {
		return ErrDuplicatePlaceholder
	}

	if len(tPlaceholders) > maxPlaceholders {
		return ErrTooManyPlaceholders
	}

	if len(tPlaceholders) < minPlaceholders {
		return ErrNotEnoughPlaceholders
	}

	if len(tPlaceholders) != len(blanks) {
		return ErrCountMismatch
	}

	for i := range tPlaceholders {
		if c := countPlaceholderEntries(tPlaceholders[i], blanks); c != 1 {
			return ErrInvalidBlanks
		}
	}

	return nil
}

func validatePlaceholder(placeholder string) error {
	if !placeholderRegexp.MatchString(placeholder) {
		return ErrInvalidPlaceholder
	}

	return nil
}

func validateAnswers(answers []option.ContentOption) error {
	if len(answers) == 0 {
		return ErrEmptyAnswers
	}

	if len(answers) > maxAnswersPerPlaceholder {
		return ErrTooManyAnswers
	}

	for i := range answers {
		if !answers[i].IsText() {
			return ErrInvalidAnswerFormat
		}
	}

	return nil
}

func countPlaceholderEntries(target string, blanks []BlankParams) int {
	count := 0
	for i := range blanks {
		if blanks[i].Placeholder == target {
			count++
		}
	}

	return count
}
