package score

type Score struct {
	value float64
}

func New(value float64) (Score, error) {
	if err := validateValue(value); err != nil {
		return Score{}, err
	}

	return Score{
		value: value,
	}, nil
}

func (s Score) Value() float64 {
	return s.value
}
