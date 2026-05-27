package education

// Education — VO сведений об образовании человека, хранит нормализованное значение.
type Education struct {
	value string
}

// New создаёт новый VO, в котором хранится нормализованное и валидированное
// значение сырой строки сведения об образовании.
func New(value string) (Education, error) {
	value = normalizeValue(value)

	if err := validateValue(value); err != nil {
		return Education{}, err
	}

	return Education{
		value: value,
	}, nil
}

// Value возвращает нормализованную строку сведений об образовании.
func (e Education) Value() string {
	return e.value
}

// IsZero возвращает true, если Education не был инициализирован через New.
func (e Education) IsZero() bool {
	return e.value == ""
}
