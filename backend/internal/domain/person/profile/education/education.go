package education

// Education is an optional normalized value description.
type Education struct {
	value string
}

// New creates an value value object.
func New(value string) (Education, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Education{}, err
	}

	return Education{
		value: value,
	}, nil
}

// Value returns the normalized value description.
func (e Education) Value() string {
	return e.value
}

// IsZero reports whether the optional value description is empty.
func (e Education) IsZero() bool {
	return e.value == ""
}
