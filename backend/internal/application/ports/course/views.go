package course

type ShortView struct {
	ID            string
	Title         string
	Published     bool
	VersionsCount int
}

type DetailedView struct {
	ID       string
	Title    string
	Versions []VersionView
}

type VersionView struct {
	ID     string
	Title  string
	Status string
}

type StudentRatingView struct {
	AccountID         string
	EnrollmentID      string
	CourseID          string
	VersionID         string
	CompletionPercent int
	CompletedItems    int
	TotalItems        int
}
