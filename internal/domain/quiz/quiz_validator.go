package quiz

import (
	"errors"
	"strings"

	"gitflic.ru/lms/internal/domain/quiz/source"
	"github.com/google/uuid"
)

var (
	ErrEmptyTitle             = errors.New("empty title")
	ErrEmptySources           = errors.New("empty sources")
	ErrNegativeAttempts       = errors.New("negative attempt limit")
	ErrNegativeTime           = errors.New("negative time limit")
	ErrNilSource              = errors.New("nil source")
	ErrCannotRemoveLastSource = errors.New("cannot remove the last source")
	ErrDuplicatedBank         = errors.New("duplicated bank")
	ErrDuplicatedSource       = errors.New("duplicated source")
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	return nil
}

func validateSources(sources []source.Source) error {
	if len(sources) == 0 {
		return ErrEmptySources
	}

	seenSources := make(map[uuid.UUID]struct{}, len(sources))
	seenBanks := make(map[uuid.UUID]struct{}, len(sources))

	for i := range sources {
		source := sources[i]
		if _, ok := seenSources[source.ID()]; ok {
			return ErrDuplicatedSource
		}

		bank := source.BankID()
		if _, ok := seenBanks[bank]; ok {
			return ErrDuplicatedBank
		}

		seenSources[source.ID()] = struct{}{}
		seenBanks[bank] = struct{}{}
	}

	return nil
}

func validateSourceToAdd(sources []source.Source, s source.Source) error {
	for i := range sources {
		if sources[i] == s {
			return ErrDuplicatedSource
		}

		if sources[i].BankID() == s.BankID() {
			return ErrDuplicatedBank
		}
	}

	return nil
}

func validateAttemptLimit(attemptLimit int) error {
	if attemptLimit < 0 {
		return ErrNegativeAttempts
	}

	return nil
}

func validateTimeLimit(timeLimit int) error {
	if timeLimit < 0 {
		return ErrNegativeTime
	}

	return nil
}
