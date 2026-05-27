package option

type Option struct {
	value string
}

func New(value string) (Option, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Option{}, err
	}

	return Option{
		value: value,
	}, nil
}

func (o Option) Value() string {
	return o.value
}

func (o Option) IsZero() bool {
	return o.value == ""
}
