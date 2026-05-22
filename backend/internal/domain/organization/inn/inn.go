package inn

// INN is a normalized and checksum-validated tax identifier.
type INN struct {
	value string
	t     Type
}

// New creates an INN value object for the supplied type.
func New(value string, t Type) (INN, error) {
	value = normalizeValue(value)

	if err := validateValue(value, t); err != nil {
		return INN{}, err
	}

	return INN{
		value: value,
		t:     t,
	}, nil
}

// Value returns the normalized INN digits.
func (inn INN) Value() string {
	return inn.value
}

// Type returns the INN type.
func (inn INN) Type() Type {
	return inn.t
}

// IsZero reports whether the INN has not been initialized.
func (inn INN) IsZero() bool {
	return inn.value == "" || !inn.t.IsValid()
}
