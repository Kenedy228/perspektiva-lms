package snils

func isSeparatorLetter(r rune) bool {
	return r == '-'
}

func isSpaceLetter(r rune) bool {
	return r == ' '
}

func isNumber(r rune) bool {
	return r >= '0' && r <= '9'
}
