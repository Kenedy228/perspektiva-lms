package education

type Education struct {
	value string
}

func New(value string) (Education, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Education{}, err
	}

	return Education{
		value: value,
	}, nil
}

func (e Education) Value() string {
	return e.value
}
