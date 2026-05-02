package name

import (
	"fmt"
)

type Name struct {
	firstName  string
	lastName   string
	middleName string
}

func New(firstName, lastName, middleName string) (Name, error) {
	firstName = normalize(firstName)
	lastName = normalize(lastName)
	middleName = normalize(middleName)

	if err := validateRequiredPart("имя", firstName); err != nil {
		return Name{}, err
	}

	if err := validateRequiredPart("фамилия", lastName); err != nil {
		return Name{}, err
	}

	if err := validateOptionalPart("отчество", middleName); err != nil {
		return Name{}, err
	}

	return Name{
		firstName:  firstName,
		lastName:   lastName,
		middleName: middleName,
	}, nil
}

func (n Name) FirstName() string {
	return n.firstName
}

func (n Name) LastName() string {
	return n.lastName
}

func (n Name) MiddleName() string {
	return n.middleName
}

func (n Name) Fullname() string {
	if n.middleName == "" {
		return fmt.Sprintf("%s %s", n.lastName, n.firstName)
	}

	return fmt.Sprintf("%s %s %s", n.lastName, n.firstName, n.middleName)
}

func (n Name) WithInitials() string {
	if n.middleName == "" {
		return fmt.Sprintf("%s %s", n.lastName, getInitial(n.firstName))
	}

	return fmt.Sprintf("%s %s%s", n.lastName, getInitial(n.firstName), getInitial(n.middleName))
}
