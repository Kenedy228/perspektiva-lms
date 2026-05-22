package snils

// SNILS объект-значение с нормализованным и провалидированным значением СНИЛС.
type SNILS struct {
	value string
}

// New создает новый объект SNILS.
func New(value string) (SNILS, error) {
	value = normalize(value)

	if err := validate(value); err != nil {
		return SNILS{}, err
	}

	return SNILS{
		value: value,
	}, nil
}

// Value возвращает нормализованное значение value.
func (s SNILS) Value() string {
	return s.value
}

// IsZero указывает, был ли инициализирован объект SNILS.
func (s SNILS) IsZero() bool {
	return s.value == ""
}
