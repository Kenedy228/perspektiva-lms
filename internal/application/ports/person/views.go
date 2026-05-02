package person

type PersonShortView struct {
	ID               string
	FullName         string
	OrganizationName string
}

type PersonDetailedView struct {
	ID               string
	FirstName        string
	LastName         string
	MiddleName       string
	Snils            string
	JobTitle         string
	Education        string
	DateOfBirth      string
	OrganizationName string
}
