package commands

// SelectableOptionInput описывает вариант ответа для вопроса с выбором.
type SelectableOptionInput struct {
	Text      string
	IsCorrect bool
}

// SequenceOptionInput описывает вариант шага для вопроса на последовательность.
type SequenceOptionInput struct {
	Text string
}

// MatchingPairInput описывает пару для вопроса на сопоставление.
type MatchingPairInput struct {
	Prompt string
	Match  string
}

// ShortVariantInput описывает допустимый вариант для короткого ответа.
type ShortVariantInput struct {
	Text string
}
