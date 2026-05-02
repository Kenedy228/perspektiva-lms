package orgname

type Name struct {
	value string
}

func New(value string) (Name, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Name{}, err
	}

	return Name{
		value: value,
	}, nil
}

func (n Name) Value() string {
	return n.value
}
