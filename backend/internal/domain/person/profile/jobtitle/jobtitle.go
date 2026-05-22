package jobtitle

// JobTitle is an optional normalized job value.
type JobTitle struct {
	value string
}

// New creates a job value value object.
func New(value string) (JobTitle, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return JobTitle{}, err
	}

	return JobTitle{
		value: value,
	}, nil
}

// Value returns the normalized job value.
func (jt JobTitle) Value() string {
	return jt.value
}

// IsZero reports whether the optional job value is empty.
func (jt JobTitle) IsZero() bool {
	return jt.value == ""
}
