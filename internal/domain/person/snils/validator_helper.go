package snils

import (
	"fmt"
	"strconv"
	"unicode"
)

func validateLength(value string) error {
	if len(value) != snilsLength {
		return fmt.Errorf("%w, детали: длина значения должна быть %d цифры", ErrInvalid, snilsLength)
	}

	return nil
}

func validateContent(value string) error {
	runes := []rune(value)

	if len(runes) == 0 {
		return fmt.Errorf("%w, детали: СНИЛС не может быть пустым", ErrInvalid)
	}

	equalDigitsCount := 0
	prevRune := unicode.ReplacementChar

	for _, r := range runes {
		if !isNumber(r) {
			return fmt.Errorf("%w, детали: СНИЛС должен содержать только цифры", ErrInvalid)
		}

		if r != prevRune {
			prevRune = r
			equalDigitsCount = 1
		} else {
			equalDigitsCount++
			if equalDigitsCount > maxEqualDigits {
				return fmt.Errorf("%w, детали: СНИЛС не может содержать более %d повторяющихся цифр подряд", ErrInvalid, maxEqualDigits)
			}
		}
	}

	return nil
}

func validateChecksum(value string) error {
	converted, err := strconv.ParseInt(value[:9], 10, 64)
	if err != nil {
		return fmt.Errorf("%w, детали: СНИЛС не является числом", ErrInvalid)
	}

	if converted <= checksumFrom {
		return nil
	}

	numPart := value[:9]
	checkPart := value[9:]

	weights := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	sum := 0

	for i, w := range weights {
		d := int(numPart[i] - '0')
		sum += d * w
	}

	switch {
	case sum < 100:
	case sum == 100 || sum == 101:
		sum = 0
	default:
		sum = sum % 101
		if sum == 100 || sum == 101 {
			sum = 0
		}
	}

	sumStr := strconv.Itoa(sum)
	if len(sumStr) == 1 {
		sumStr = "0" + sumStr
	}

	if sumStr != checkPart {
		return fmt.Errorf("%w: некорректная контрольная сумма", ErrInvalid)
	}

	return nil
}
