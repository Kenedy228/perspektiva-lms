package login

type Login struct {
	value string
}

func New(value string) (Login, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Login{}, err
	}

	return Login{
		value: value,
	}, nil
}

func (l Login) Value() string {
	return l.value
}
