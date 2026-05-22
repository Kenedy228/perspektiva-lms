package name

// Name объект-значение имени человека в кириллице.
type Name struct {
	firstName  string
	lastName   string
	middleName string
}

// New создает новое имя человека.
func New(firstName, lastName, middleName string) (Name, error) {
	firstName = normalizeFirstName(firstName)
	lastName = normalizeLastName(lastName)
	middleName = normalizeMiddleName(middleName)

	if err := validateFirstName(firstName); err != nil {
		return Name{}, err
	}
	if err := validateLastName(lastName); err != nil {
		return Name{}, err
	}
	if err := validateMiddleName(middleName); err != nil {
		return Name{}, err
	}

	return Name{
		firstName:  firstName,
		lastName:   lastName,
		middleName: middleName,
	}, nil
}

// FirstName возвращает имя.
func (n Name) FirstName() string {
	return n.firstName
}

// LastName возвращает фамилию.
func (n Name) LastName() string {
	return n.lastName
}

// MiddleName возвращает отчество.
func (n Name) MiddleName() string {
	return n.middleName
}

// IsZero сигнализирует, был ли инициализирован объект Name.
func (n Name) IsZero() bool {
	return n.firstName == "" || n.lastName == ""
}
