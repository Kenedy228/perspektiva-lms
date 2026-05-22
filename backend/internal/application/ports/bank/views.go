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
}
