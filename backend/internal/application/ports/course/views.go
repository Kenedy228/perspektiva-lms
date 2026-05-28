package course

type ShortView struct {
	ID          string
	Title       string
	Published   bool
	BlocksCount int
}

type DetailedView struct {
	ID    string
	Title string
}

type StudentRatingView struct {
	AccountID         string
	EnrollmentID      string
	CourseID          string
	CompletionPercent int
	CompletedItems    int
	TotalItems        int
}
