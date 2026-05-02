package title

type Title struct {
	value string
}

func New(value string) (Title, error) {
	if err := validateValue(value); err != nil {
		return Title{}, err
	}

	return Title{
		value: value,
	}, nil
}

func (t Title) Value() string {
	return t.value
}
