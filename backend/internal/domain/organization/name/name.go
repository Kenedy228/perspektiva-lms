package name

// Name is a normalized organization name.
type Name struct {
	value string
}

// New creates an organization name value object.
func New(value string) (Name, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Name{}, err
	}

	return Name{
		value: value,
	}, nil
}

// Value returns the normalized organization name.
func (n Name) Value() string {
	return n.value
}

// IsZero reports whether the name has not been initialized.
func (n Name) IsZero() bool {
	return n.value == ""
}
