package bank

import (
	"encoding/binary"
	"errors"
	"slices"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrEmptyTitle        = errors.New("title can't be empty")
	ErrNilQuestion       = errors.New("nil question")
	ErrQuestionDuplicate = errors.New("question already present in bank")
)

func validateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}

	return nil
}

func validateQuestion(questionID uuid.UUID) error {
	if questionID == uuid.Nil {
		return ErrNilQuestion
	}

	return nil
}

func validateQuestionsForAdding(current, add []uuid.UUID) error {
	const wordsCount = 512
	const totalBits = wordsCount * 64
	bitset := [wordsCount]uint64{}

	for i := range add {
		if err := validateQuestion(add[i]); err != nil {
			return err
		}

		hash := binary.LittleEndian.Uint32(add[i][12:])

		absoluteBitIndex := hash % totalBits
		wordIndex := absoluteBitIndex / 64
		bitOffset := (absoluteBitIndex % 64)
		bitMask := uint64(1 << bitOffset)

		if (bitset[wordIndex] & bitMask) != 0 {
			if slices.Contains(add[:i], add[i]) {
				return ErrQuestionDuplicate
			}
		}

		if slices.Contains(current, add[i]) {
			return ErrQuestionDuplicate
		}

		bitset[wordIndex] |= bitMask
	}

	return nil
}
