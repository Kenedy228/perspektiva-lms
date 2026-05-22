package variant

type Variant struct {
	value string
}

func New(value string) (Variant, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Variant{}, err
	}

	return Variant{
		value: value,
	}, nil
}

func (v Variant) Value() string {
	return v.value
}

func (v Variant) IsZero() bool {
	return v.value == ""
}
