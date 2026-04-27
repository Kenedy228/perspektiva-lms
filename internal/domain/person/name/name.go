package name

import (
	"fmt"
	"unicode/utf8"
)

type Name struct {
	firstname  string
	lastname   string
	middlename string
}

func New(params Params) (Name, error) {
	firstname := normalize(params.Firstname)
	lastname := normalize(params.Lastname)
	middlename := normalize(params.Middlename)

	if err := validateRequiredPart("имя", firstname); err != nil {
		return Name{}, err
	}

	if err := validateRequiredPart("фамилия", lastname); err != nil {
		return Name{}, err
	}

	if err := validateOptionalPart("отчество", middlename); err != nil {
		return Name{}, err
	}

	return Name{
		firstname:  firstname,
		lastname:   lastname,
		middlename: middlename,
	}, nil
}

func (n Name) Firstname() string {
	return n.firstname
}

func (n Name) Lastname() string {
	return n.lastname
}

func (n Name) Middlename() string {
	return n.middlename
}

func (n Name) Fullname() string {
	if n.middlename == "" {
		return fmt.Sprintf("%s %s", n.lastname, n.firstname)
	}

	return fmt.Sprintf("%s %s %s", n.lastname, n.firstname, n.middlename)
}

func (n Name) WithInitials() string {
	firstnameInitial, _ := utf8.DecodeRuneInString(n.firstname)

	if n.middlename == "" {
		return fmt.Sprintf("%s %c.", n.lastname, firstnameInitial)
	}

	middlenameInitial, _ := utf8.DecodeRuneInString(n.middlename)
	return fmt.Sprintf("%s %c.%c.", n.lastname, firstnameInitial, middlenameInitial)
}
