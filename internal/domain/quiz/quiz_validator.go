package quiz

import (
	"errors"
	"strings"

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
)

func validateParams() error {

}

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	return nil
}

func validateSources(sources []Source) error {
	if len(sources) == 0 {
		return ErrEmptySources
	}

	seenBanks := make(map[uuid.UUID]struct{}, len(sources))

	for i := range sources {
		bank := sources[i].BankID()
		if _, ok := seenBanks[bank]; ok {
			return ErrDuplicatedBank
		}

		seenBanks[bank] = struct{}{}
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
