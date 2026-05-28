package bank

type ShortView struct {
	ID             string
	Title          string
	QuestionsCount int
}

type DetailedView struct {
	ID          string
	Title       string
	QuestionIDs []string
	Questions   []QuestionView
}

type QuestionView struct {
	ID                string
	Type              string
	Title             string
	SelectableOptions []SelectableOptionView
	SequenceOptions   []SequenceOptionView
	MatchingPairs     []MatchingPairView
	ShortVariants     []ShortVariantView
}

type SelectableOptionView struct {
	ID        string
	Value     string
	IsCorrect bool
}

type SequenceOptionView struct {
	Value string
}

type MatchingPairView struct {
	PromptID   string
	PromptText string
	MatchID    string
	MatchText  string
}

type ShortVariantView struct {
	Value string
}
