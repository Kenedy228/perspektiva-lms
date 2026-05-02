package jobtitle

type JobTitle struct {
	title string
}

func New(title string) (JobTitle, error) {
	title = normalizeJobTitle(title)

	if err := validateJobTitle(title); err != nil {
		return JobTitle{}, err
	}

	return JobTitle{
		title: title,
	}, nil
}

func (jt JobTitle) Title() string {
	return jt.title
}
