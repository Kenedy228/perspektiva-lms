package commands

type AttachmentInput struct {
	MediaType string
	FileName  string
	SizeBytes int64
}

type SelectableOptionInput struct {
	Text      string
	IsCorrect bool
}

type SequenceOptionInput struct {
	Text string
}

type MatchingPairInput struct {
	Prompt string
	Match  string
}

type TypedBlankInput struct {
	Placeholder string
	Variants    []string
}

type ShortVariantInput struct {
	Text string
}
