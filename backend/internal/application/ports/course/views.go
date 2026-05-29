package course

type ShortView struct {
	ID          string
	Title       string
	Published   bool
	BlocksCount int
}

type DetailedView struct {
	ID     string
	Title  string
	Blocks []BlockView
}

type BlockView struct {
	ID       string
	Title    string
	Elements []ElementView
}

type ElementView struct {
	ID             string
	Title          string
	Type           string
	CompletionMode string
	QuizID         string
}

type StudentRatingView struct {
	AccountID         string
	EnrollmentID      string
	CourseID          string
	CompletionPercent int
	CompletedItems    int
	TotalItems        int
}
