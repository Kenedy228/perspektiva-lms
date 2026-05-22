package login

// Login объект-значение над логином аккаунта.
type Login struct {
	value string
}

// New создает новый объект Login.
func New(value string) (Login, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Login{}, err
	}

	return Login{
		value: value,
	}, nil
}

// Value возвращает значение логина.
func (l Login) Value() string {
	return l.value
}

// IsZero сигнализирует, была ли инициализация объекта Login.
func (l Login) IsZero() bool {
	return l.value == ""
}
