package inn

import (
	"unicode"
)

func getDigits(s string) ([]int, bool) {
	digits := make([]int, 0, len(s))

	for _, r := range s {
		if !unicode.IsDigit(r) {
			return nil, false
		}

		digits = append(digits, int(r-'0'))
	}

	return digits, true
}
