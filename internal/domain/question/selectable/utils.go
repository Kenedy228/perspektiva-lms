package selectable

func countCorrect(options []Option) int {
	count := 0

	for i := range options {
		if options[i].IsCorrect() {
			count++
		}
	}

	return count
}
