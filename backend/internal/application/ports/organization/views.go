package organization

type OrganizationShortView struct {
	ID               string
	OrganizationName string
}

type OrganizationDetailedView struct {
	ID               string
	OrganizationName string
	INN              string
	INNTitle         string
}
