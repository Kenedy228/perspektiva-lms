package snils

func normalize(s string) string {
	runes := []rune(s)
	normalized := make([]rune, len(s))
	idx := 0

	for _, r := range runes {
		if !isSeparatorLetter(r) && !isSpaceLetter(r) {
			normalized[idx] = r
			idx++
		}
	}

	return string(normalized[0:idx])
}
